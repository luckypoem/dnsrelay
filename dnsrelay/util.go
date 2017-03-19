package dnsrelay

import (
	"github.com/miekg/dns"
	"net"
)

func GetIPFromMsg(r *dns.Msg) (addrs []net.IP, err error) {
	for _, a := range r.Answer {
		switch ta := a.(type) {
		case *dns.A:
			addrs = append(addrs, ta.A)
		case *dns.AAAA:
			addrs = append(addrs, ta.AAAA)
		}
	}

	return
}
