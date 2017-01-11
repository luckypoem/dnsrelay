package dnsrelay

import (
	"net"
	"sync"
	"github.com/miekg/dns"
	"github.com/FTwOoO/go-logger"
	"github.com/FTwOoO/vpncore/net/geoip"
	"github.com/FTwOoO/vpncore/net/rule"
	"github.com/FTwOoO/vpncore/net/addr"

	"time"
	"errors"
	"fmt"
)

const (
	dnsDefaultPort = 53
	dnsDefaultTtl = 600
	dnsDefaultPacketSize = 4096
	dnsDefaultReadTimeout = 5
	dnsDefaultWriteTimeout = 5
)

type dnsServerInternal interface {
	ListenAndServe() error
}

type DNSServer struct {
	config    *Config
	geoReader *geoip.Reader
	logger    *logger.Logger
	cache     *MemoryCache
	client    *dns.Client
}

func CreateNormalDnsServer(addr string) (dnsServerInternal, error) {
	server := &dns.Server{
		Net:          "udp",
		Addr:         addr,
		UDPSize:      dnsDefaultPacketSize,
		ReadTimeout:  time.Duration(dnsDefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(dnsDefaultWriteTimeout) * time.Second,
	}

	return server, nil
}

func NewDNSServer(config *Config, customDnsServer dnsServerInternal) (ds *DNSServer, err error) {
	reader, err := geoip.Open(config.GeoIPDBPath)
	if err != nil {
		return nil, err
	}

	var cache MemoryCache
	switch config.DNSCache.Backend {
	case "memory":
		cache = MemoryCache{
			Backend:  make(map[string]DomainRecord, config.DNSCache.Maxcount),
			DefaultTtl:   time.Duration(config.DNSCache.Expire),
			Maxcount: config.DNSCache.Maxcount,
		}
		cache.Serve()
	default:
		return nil, errors.New("Cache backend dont support!")
	}

	client := &dns.Client{
		Net:          "udp",
		UDPSize:      dnsDefaultPacketSize,
		ReadTimeout:  time.Duration(dnsDefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(dnsDefaultWriteTimeout) * time.Second,
	}

	mlogger, err := logger.NewLogger(config.LogConfig.LogFile, config.LogConfig.LogLevel)
	if err != nil {
		return
	}

	ds = &DNSServer{
		config:     config,
		geoReader:  reader,
		logger:     mlogger,
		cache:      &cache,
		client:     client,
	}

	var inDnsServ dnsServerInternal
	if customDnsServer == nil {
		inDnsServ, _ = CreateNormalDnsServer(fmt.Sprintf("%s:%d", config.ADDR, config.PORT))
		if inDs, ok := inDnsServ.(*dns.Server); ok {
			inDs.Handler = ds
		}
	}
	go inDnsServ.ListenAndServe()
	return
}


// Main callback for miekg/dns. Collects information about the query,
// constructs a response, and returns it to the connector.
func (ds *DNSServer) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {

	if len(req.Question) == 0 {
		dns.HandleFailed(w, req)
		return
	}

	var resp *dns.Msg
	hitCache := false

	question := req.Question[0]
	if question.Qclass == dns.ClassINET &&
		(question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA) {

		if cacheResp, err := ds.cache.Get(question.Name); err == nil {
			ds.logger.Debugf("%s hit cache", question.Name)
			// dont change cache object, copy it
			newResp := *cacheResp
			newResp.Id = req.Id
			resp = &newResp
			hitCache = true
		} else if newResp, ok := ds.config.Hosts.Get(req); ok {
			ds.logger.Debugf("%s found in hosts file", question.Name)
			resp = newResp
		}
	}

	if resp == nil {
		group := ds.config.DomainRules.FindGroup(question.Name)
		if group == rule.REJECT_GROUP {
			ds.logger.Debugf("Reject %s!", question.Name)
			resp = nil
		} else if group != "" {
			resp = ds.sendRequest(req, []string{group})
		} else {
			resp = ds.sendRequest(req, ds.config.DefaultGroups)
		}
	}

	if resp != nil {
		w.WriteMsg(resp)

		if question.Qclass == dns.ClassINET &&
			(question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA) &&
			len(resp.Answer) > 0 {

			if hitCache == true {
				ds.logger.Debugf("No need to insert %s into cache", question.Name)

			} else if err := ds.cache.Set(question.Name, resp, time.Duration(resp.Answer[0].Header().Ttl)); err != nil {
				ds.logger.Warningf("Set %s cache failed: %s", question.Name, err.Error())
			} else {
				ds.logger.Debugf("Insert %s into cache", question.Name)
			}
		}
	} else {
		dns.HandleFailed(w, req)
	}
}

type DNSResult struct {
	Group    string
	DnsIp    net.IP

	Response *dns.Msg
	Rtt      time.Duration
	Err      error
}

func (ds *DNSServer) sendRequest(req *dns.Msg, dnsgroups []string) (resp *dns.Msg) {

	chanLen := 0
	for _, group := range dnsgroups {
		chanLen += len(ds.config.DNSGroups[group])
	}

	// the chan is big enough to hold all results
	results := make(chan DNSResult, chanLen)
	ds.sendDNSRequestsAsync(req, results, dnsgroups)

	WaitingDNSResponse:
	for result := range results {
		if result.Err != nil {
			ds.logger.Errorf("Error from group[%s] DNS[%s]: ===>\n %v \n<===\n", result.Group, result.DnsIp.String(), result.Err)
			continue
		} else {
			ds.logger.Debugf("Result from group[%s] DNS[%s]: ===>\n %v \n<===\n", result.Group, result.DnsIp.String(), result.Response)
		}

		if result.Response.Rcode == dns.RcodeServerFailure {
			ds.logger.Errorf("Resolve on group [%s:%s] failed: code %d", result.Group, result.DnsIp.String(), result.Response.Rcode)
			continue
		}

		if len(result.Response.Answer) < 1 {
			ds.logger.Debugf("0 answer response from  %s", result.DnsIp.String())
			continue
		}

		switch result.Response.Answer[0].(type) {
		case *dns.A:
			aRecord, _ := result.Response.Answer[0].(*dns.A)
			resultIp := aRecord.A

			if ds.isIpOK(result.Group, resultIp) {
				resp = result.Response
				break WaitingDNSResponse
			}
		default:
			resp = result.Response
			break WaitingDNSResponse
		}
	}

	ds.logger.Debugf("response for request:%v\n", resp)
	return resp
}

func (ds *DNSServer) isIpOK(dnsGroup string, resultIp net.IP) bool {
	if ds.config.IPFilter.FindIP(resultIp) {
		ds.logger.Infof("block ip %v", resultIp.String())
		return false
	}

	if ds.config.FuckGFW == false {
		return true
	}

	isCN, err := ds.geoReader.IsChineseIP(resultIp)

	if err != nil {
		ds.logger.Errorf("cant reconize the localtion of ip:%s", resultIp.String())
		if dnsGroup != rule.CN_GROUP {
			return true
		} else {
			return false
		}
	}

	if dnsGroup == rule.CN_GROUP  && isCN {
		ds.logger.Debugf("DNS result of CN IP[%s] from CN DNS server can be trusted!", resultIp)
		return true
	} else if dnsGroup != rule.CN_GROUP  && !isCN {
		ds.logger.Debugf("DNS result of NO CN IP[%s] from NO CN DNS server can be trusted!", resultIp)
		return true
	}

	return false
}

func (ds *DNSServer) sendDNSRequestsAsync(req *dns.Msg, results chan <- DNSResult, dnsgroups []string) {
	var wg sync.WaitGroup

	for _, group := range dnsgroups {
		dnsL := ds.config.DNSGroups[group]
		wg.Add(len(dnsL))

		for _, dnsAddr := range dnsL {
			go func(group string, dnsAddr addr.DNSAddresss) {
				defer wg.Done()

				//c := &dns.Client{Net: "udp", Timeout:10 * time.Second}
				resp, rtt, err := ds.client.Exchange(req, dnsAddr.String())
				results <- DNSResult{Response:resp, Rtt: rtt, Err: err, Group:group, DnsIp:dnsAddr.Ip}

			}(group, dnsAddr)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

}