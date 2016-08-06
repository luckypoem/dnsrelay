package dnsrelay

import (
	"testing"
	"net"
	"fmt"
		"math/rand"

)

func TestDomainRules(t *testing.T) {
	rules := DomainRules{Rules: *new([]Rule)}
	rules.AddRule(Rule{Group:CN_GROUP, MatchType:MATCH_TYPE_DOMAIN_SUFFIX, Value:"baidu.com"})
	rules.AddRule(Rule{Group:CN_GROUP, MatchType:MATCH_TYPE_DOMAIN_KEYWORD, Value:"xiaomi"})
	rules.AddRule(Rule{Group:FG_GROUP, MatchType:MATCH_TYPE_DOMAIN, Value:"www.google.com"})

	group := rules.FindGroup("baidu.com")
	if group != CN_GROUP {
		t.Fail()
	}

	group = rules.FindGroup("xiaomi.com")
	if group != CN_GROUP {
		t.Fail()
	}

	group = rules.FindGroup("www.google.com")
	if group != FG_GROUP {
		t.Fail()
	}

	group = rules.FindGroup("google.com")
	if group != "" {
		t.Fail()
	}

}

func BenchmarkNewBlackIP(b *testing.B) {

	r := rand.New(rand.NewSource(0))

	ips := []net.IP{}
	for i := 0; i < b.N; i++ {
		ip := randomIPv4Address(b, r)
		ips = append(ips, ip)
	}


	iprule := NewBlackIP(ips)

	for i := 0; i < b.N; i++ {
		ip := ips[i]
		f := iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)

		ip = randomIPv4Address(b, r)
		f = iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)
	}


}