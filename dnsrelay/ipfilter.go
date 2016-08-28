package dnsrelay

import (
	"sort"
	"net"
)

type IPFilter struct {
	Ip  IPList
	Net IPNetList
}

func NewIPBlocker(ips IPList) (*IPFilter) {
	sort.Sort(ips)
	return &IPFilter{Ip:ips}
}

func (self *IPFilter) FindIP(ip net.IP) bool {
	return self.Ip.Contains(ip) || self.Net.Contains(ip)
}