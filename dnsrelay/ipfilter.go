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
	i := sort.Search(len(self.Ip), func(i int) bool {
		// TODO: support IPv6
		return IPToInt(self.Ip[i]) >= IPToInt(ip)
	})
	return i < len(self.Ip) && self.Ip[i].Equal(ip)
}