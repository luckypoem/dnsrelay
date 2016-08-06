package dnsrelay

import (
	"net"
	"sync"

	"github.com/miekg/dns"
	"fmt"
)

type DNSServer struct {
	aRecords map[string]net.IP // FQDN -> IP
	aMutex   sync.RWMutex      // mutex for A record operations

	config   Config
	reader   *Reader
}

// Create a new DNS server. Domain is an unqualified domain that will be used
// as the TLD.
func NewDNSServer(config Config) (*DNSServer, error) {
	reader, err := Open(config.GeoIPDBPath)
	if err != nil {
		return nil, err
	}

	return &DNSServer{
		aRecords:   map[string]net.IP{},
		aMutex:     sync.RWMutex{},
		config:     config,
		reader: reader,
	}, nil
}

// Listen for DNS requests. listenSpec is a dotted-quad + port, e.g.,
// 127.0.0.1:53. This function blocks and only returns when the DNS service is
// no longer functioning.
func (ds *DNSServer) Listen(listenSpec string) error {
	return dns.ListenAndServe(listenSpec, "udp", ds)
}



// Receives a FQDN; looks up and supplies the A record.
func (ds *DNSServer) GetA(fqdn string) *dns.A {
	ds.aMutex.RLock()
	defer ds.aMutex.RUnlock()
	val, ok := ds.aRecords[fqdn]

	if ok {
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   fqdn,
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

func (ds *DNSServer) route(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0  {
		dns.HandleFailed(w, req)
		return
	}
	group := ds.config.DomainRules.FindGroup(req.Question[0].Name)

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
	err      error
}

func (ds *DNSServer) proxy(w dns.ResponseWriter, req *dns.Msg, dnsgroups []string) {

	counter := 0

	for _, group := range dnsgroups {
		counter += len(ds.config.DNSGroups[group])
	}

	results := make(chan DNSResult, counter)
	ds.sendDNSRequestsAsync(req, results, dnsgroups)

	var response *dns.Msg

	for result := range results {

		counter -= 1
		if counter == 0 {
			break
		}

		if result.err != nil {
			continue
		}

		aRecord, ok := result.Response.Answer[0].(*dns.A)
		if !ok {
			fmt.Printf("Not a A record return %v", aRecord)
			continue
		}

		result_ip := aRecord.A

		if !ds.config.BlackIP.FindIP(result_ip) {

			if result.group == CN_GROUP {
				if is, err :=ds.reader.IsChineseIP(result_ip); err != nil && is{
					response = result.Response
					break
				}
			} else {
				if is, err :=ds.reader.IsChineseIP(result_ip); err != nil && !is {
					response = result.Response
					break
				}
			}
		}

	}

	close(results)

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

	for _, group := range dnsgroups {
		dns_ips := ds.config.DNSGroups[group]
		for _, dns_ip := range dns_ips {
			go func() {
				// err.. is there a async func for dns.Client?
				c := &dns.Client{Net: "udp"}
				// 2 seconds to timeout
				resp, _, err := c.Exchange(req, dns_ip.String())
				if err == nil {
					results <- DNSResult{Response:resp, group:group}
				} else {
					results <- DNSResult{err:err}
				}

				return
			}()
		}
	}

}