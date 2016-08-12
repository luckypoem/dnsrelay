package dnsrelay

import "github.com/miekg/dns"

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s) - 1]
	}
	return s
}

