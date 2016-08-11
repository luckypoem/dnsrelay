package dnsrelay

import (
	"net"
	"sort"
	"encoding/binary"
)

func IPToInt(ip net.IP) uint32 {
	// net.ParseIP will return 16 bytes for IPv4,
	// but we cant stop user creating 4 bytes for IPv4 bytes using net.IP{N,N,N,N}
	//
	if len(ip) == net.IPv4len {
		return binary.BigEndian.Uint32(ip)
	} else if len(ip) == net.IPv6len {
		return binary.BigEndian.Uint32(ip[12:])
	}

	return 0
}

type IPList []net.IP

func (a IPList) Len() int {
	return len(a)
}
func (a IPList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a IPList) Less(i, j int) bool {
	return IPToInt(a[i]) < IPToInt(a[j])
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
		// TODO: support IPv6
		return IPToInt(self.ips[i]) >= IPToInt(ip)
	})
	return i < len(self.ips) && self.ips[i].Equal(ip)
}
