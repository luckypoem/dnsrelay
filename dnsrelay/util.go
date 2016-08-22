package dnsrelay

import (
	"github.com/miekg/dns"
	"github.com/miekg/dns/dnsutil"
)

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return dnsutil.TrimDomainName(s, ".")
	}
	return s
}

