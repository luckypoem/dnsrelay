package dnsrelay

import (
	"net"
	"sort"
)

func IPToU32(ip net.IP) uint32 {
	return (uint32)(ip[0]) << 24 | (uint32)(ip[1]) << 16 | (uint32)(ip[2]) << 8 | (uint32)(ip[3])
}

type IPList []net.IP

func (a IPList) Len() int {
	return len(a)
}
func (a IPList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a IPList) Less(i, j int) bool {
	return IPToU32(a[i]) < IPToU32(a[j])
}

type IPBlocker struct {
	ips IPList
}

func NewIPBlocker(ips IPList) (*IPBlocker) {
	sort.Sort(ips)
	return &IPBlocker{ips:ips}
}

func (self *IPBlocker) FindIP(ip net.IP) bool {
	i := sort.Search(len(self.ips), func(i int) bool {
		return IPToU32(self.ips[i]) >= IPToU32(ip)
	})
	return i < len(self.ips) && self.ips[i].Equal(ip)
}
