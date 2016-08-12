package dnsrelay

import (
	"github.com/miekg/dns"
)

type Host struct {
	Name string `toml:"name"`
	Ip   IPList `toml:"ip"`
}

type Hosts []Host

func (self Hosts) Get(req *dns.Msg) (*dns.Msg, bool) {
	question := req.Question[0]

	var ips IPList
	queryHost := UnFqdn(question.Name)

	for _, host := range self {
		if host.Name == queryHost {
			ips = host.Ip
			break
		}
	}

	if ips == nil {
		return nil, false
	}

	resp := new(dns.Msg)
	resp.SetReply(req)

	switch question.Qtype {
	case dns.TypeA:
		rr_header := dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    1,
		}
		for _, ip := range ips {
			ip = ip.To4()
			if ip != nil {
				a := &dns.A{rr_header, ip}
				resp.Answer = append(resp.Answer, a)
			}
		}
	case dns.TypeAAAA:
		rr_header := dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    1,
		}
		for _, ip := range ips {
			ip = ip.To16()
			if ip != nil {
				aaaa := &dns.AAAA{rr_header, ip}
				resp.Answer = append(resp.Answer, aaaa)
			}
		}

	default:
		return nil, false
	}

	return resp, true
}



