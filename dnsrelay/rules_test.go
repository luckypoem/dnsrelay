package dnsrelay

import (
	"testing"
	"net"
	"fmt"
	"math/rand"

)

func TestDomainRules(t *testing.T) {
	rule := DomainRule{Group:CN_GROUP,
		MatchType:SCHEME_DOMAIN_SUFFIX,
		Values:[]string{"baidu.com", "xiaomi.com"},
	}

	r := rule.Match("baidu.com")
	if r != true {
		t.Fail()
	}

	r = rule.Match("xiaomi.com")
	if r != true {
		t.Fail()
	}

	r = rule.Match("www.google.com")
	if r != false {
		t.Fail()
	}

	r = rule.Match("google.com")
	if r != false {
		t.Fail()
	}

	rules := DomainRules{rule}

	g := rules.findGroup("baidu.com")
	if g != CN_GROUP {
		t.Fatal("rules match fail!")
	}

}

func BenchmarkNewBlackIP(b *testing.B) {

	r := rand.New(rand.NewSource(0))

	ips := []net.IP{}
	for i := 0; i < b.N; i++ {
		ip := randomIPv4Address(b, r)
		ips = append(ips, ip)
	}

	iprule := NewIPBlocker(ips)

	for i := 0; i < b.N; i++ {
		ip := ips[i]
		f := iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)

		ip = randomIPv4Address(b, r)
		f = iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)
	}

}