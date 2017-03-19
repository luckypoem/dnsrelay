package dnsrelay

import (
	"testing"

	"github.com/miekg/dns"
	"fmt"
)

const (
	nameserver = "127.0.0.1:53"
	domain = "www.sina.com.cn"
)

func BenchmarkDig(b *testing.B) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	c := new(dns.Client)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Exchange(m, nameserver)
	}

}

func callbackFunc(msg *dns.Msg) error {
	ips, err := GetIPFromMsg(msg)
	if err != nil {
		return err
	}

	fmt.Print("result ip:", ips)
	return nil
}

func TestDnsServer(t *testing.T) {
	ds, err := NewDNSServer(nil, true)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	ds.QueryIPv4("baidu.com", callbackFunc)

}