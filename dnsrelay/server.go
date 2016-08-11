package dnsrelay

import (
	"net"
	"sync"
	"github.com/miekg/dns"
	"fmt"
	"time"
)

type DNSServer struct {
	aRecords map[string]net.IP
	aMutex   sync.RWMutex // mutex for A record operations

	config   *Config
	reader   *Reader
	logger   *Logger
}

// Create a new DNS server. Domain is an unqualified domain that will be used
// as the TLD.
func NewDNSServer(config *Config) (*DNSServer, error) {
	reader, err := Open(config.GeoIPDBPath)
	if err != nil {
		return nil, err
	}

	ds := DNSServer{
		aRecords:   map[string]net.IP{},
		aMutex:     sync.RWMutex{},
		config:     config,
		reader:     reader,
		logger:     config.Logger,
	}

	ds.SetupHosts(config)

	return &ds, nil
}

// Listen for DNS requests. listenSpec is a dotted-quad + port, e.g.,
// 127.0.0.1:53. This function blocks and only returns when the DNS service is
// no longer functioning.
func (ds *DNSServer) Listen(listenSpec string) error {
	ds.logger.Infof("Listen on %s ...", listenSpec)
	return dns.ListenAndServe(listenSpec, "udp", ds)
}

// looks up and supplies the A record.
func (ds *DNSServer) GetA(domain string) *dns.A {
	ds.aMutex.RLock()
	defer ds.aMutex.RUnlock()
	val, ok := ds.aRecords[domain]

	if ok {
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   domain,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				// 0 TTL results in UB for DNS resolvers and generally causes problems.
				Ttl: 1,
			},
			A: val,
		}
	}

	return nil
}

// Sets a host to an IP. Note that this is not the FQDN, but a hostname.
func (ds *DNSServer) SetA(host string, ip net.IP) {
	ds.aMutex.Lock()
	ds.aRecords[host] = ip
	ds.aMutex.Unlock()
}

// Deletes a host. Note that this is not the FQDN, but a hostname.
func (ds *DNSServer) DeleteA(host string) {
	ds.aMutex.Lock()
	delete(ds.aRecords, host)
	ds.aMutex.Unlock()
}

func (ds *DNSServer) CleanCache() {
	ds.aRecords = map[string]net.IP{}
}

func (ds *DNSServer) SetupHosts(config *Config) {
	for _, host := range config.Hosts {
		ds.SetA(host.Name, net.ParseIP(host.Ip))
	}
}

// Main callback for miekg/dns. Collects information about the query,
// constructs a response, and returns it to the connector.
func (ds *DNSServer) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {

	m := &dns.Msg{}
	m.SetReply(req)

	answers := []dns.RR{}


	// check if cache already have the answers
	for _, question := range req.Question {
		// nil records == not found
		switch question.Qtype {
		case dns.TypeA:
			a := ds.GetA(question.Name)
			if a != nil {
				answers = append(answers, a)
			}
		}
	}

	if len(answers) > 0 {
		// Without these the glibc resolver gets very angry.
		m.Authoritative = true
		m.RecursionAvailable = true
		m.Answer = answers
		w.WriteMsg(m)
	} else {
		ds.route(w, req)
	}
}

func (self *DNSServer) findGroup(domain string) string {
	for _, rule := range self.config.Rules {
		if rule.Match(domain) {
			return rule.Group
		}
	}
	return ""
}

func (ds *DNSServer) route(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		dns.HandleFailed(w, req)
		return
	}
	group := ds.findGroup(req.Question[0].Name)

	if group != "" {
		ds.proxy(w, req, []string{group})
	} else if group == REJECT_GROUP {
		dns.HandleFailed(w, req)
		return
	} else {
		ds.proxy(w, req, ds.config.DefaultGroups)
	}
}

type DNSResult struct {
	Response *dns.Msg
	group    string
	dnsIp    net.IP
	err      error
}

func (ds *DNSServer) proxy(w dns.ResponseWriter, req *dns.Msg, dnsgroups []string) {

	chanLen := 0
	for _, group := range dnsgroups {
		chanLen += len(ds.config.DNSGroups[group])
	}

	// the chan is big enough to hold all results
	results := make(chan DNSResult, chanLen)
	ds.sendDNSRequestsAsync(req, results, dnsgroups)

	var response *dns.Msg

	for result := range results {
		if result.err != nil {
			ds.logger.Debugf("Error from group(%s) DNS(%s): ===>\n %v \n<===\n", result.group, result.dnsIp.String(), result.err)
			continue
		} else {
			ds.logger.Debugf("Result from group(%s) DNS(%s): ===>\n %v \n<===\n", result.group, result.dnsIp.String(), result.Response)

		}

		aRecord, ok := result.Response.Answer[0].(*dns.A)
		if !ok {
			ds.logger.Infof("Not a A record return %v", aRecord)
			continue
		}

		resultIp := aRecord.A

		if ds.config.IPBlocker.FindIP(resultIp) {
			ds.logger.Infof("block ip %v", resultIp.String())
			continue
		}

		isCN, err := ds.reader.IsChineseIP(resultIp)

		if err != nil {
			ds.logger.Errorf("cant reconize the localtion of ip:%s", resultIp.String())
			if result.group != CN_GROUP {
				response = result.Response
				break
			} else {
				continue
			}
		}

		if result.group == CN_GROUP  && isCN {
			response = result.Response
			break
		} else if result.group != CN_GROUP  && !isCN {
			response = result.Response
			break
		}
	}

	if response == nil {
		dns.HandleFailed(w, req)
		// If we have no answers, that means we found nothing or didn't get a query
		// we can reply to. Reply with no answers so we ensure the query moves on to
		// the next server.
		//if len(answers) == 0 {
		//	m.SetRcode(r, dns.RcodeSuccess)
		//	w.WriteMsg(m)
		//	return
		//}

	} else {
		w.WriteMsg(response)
	}
}

func (ds *DNSServer) sendDNSRequestsAsync(req *dns.Msg, results chan <- DNSResult, dnsgroups []string) {
	var wg sync.WaitGroup

	for _, group := range dnsgroups {
		dns_ips := ds.config.DNSGroups[group]
		wg.Add(len(dns_ips))

		for _, dns_ip := range dns_ips {
			go func(group, dns_ip string) {
				defer wg.Done()

				// err.. is there a async func for dns.Client?
				c := &dns.Client{Net: "udp", Timeout:10 * time.Second}
				// 2 seconds to timeout
				resp, _, err := c.Exchange(req, dns_ip + ":53")
				if err == nil {
					results <- DNSResult{Response:resp, group:group, dnsIp:net.ParseIP(dns_ip)}
				} else {
					fmt.Println(err)
					results <- DNSResult{err:err, group:group, dnsIp:net.ParseIP(dns_ip)}
				}
			}(group, dns_ip)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

}