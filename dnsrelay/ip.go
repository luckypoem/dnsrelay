package dnsrelay

import (
	"net"
	"encoding/binary"
	"strings"
	"fmt"
	"strconv"
	"sort"
)

type DNSAddresss struct {
	Ip   net.IP
	Port int
}

func (self *DNSAddresss) UnmarshalTOML(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	spliter := ":"

	if strings.Contains(s, spliter) {
		arr := strings.Split(s, spliter)
		if len(arr) != 2 {
			return fmt.Errorf("Bad format for DNS:%s", s)
		}
		ip, port := arr[0], arr[1]
		self.Ip = net.ParseIP(ip)
		self.Port, err = strconv.Atoi(port)
		if err != nil {
			return err
		}
	} else {
		self.Ip = net.ParseIP(s)
		if self.Ip == nil {
			return fmt.Errorf("Bad format for DNS:%s", s)
		}

		self.Port = 53
	}
	return nil
}

func (self *DNSAddresss) String() string {
	return fmt.Sprintf("%s:%d", self.Ip.String(), self.Port)
}

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

func (a IPList) Sort() {
	sort.Sort(a)
}

func (self IPList) Contains(ip net.IP) bool {
	// TODO: support IPv6

	i := sort.Search(len(self), func(i int) bool {
		return IPToInt(self[i]) >= IPToInt(ip)
	})
	return i < len(self) && self[i].Equal(ip)
}

func (self *IPList) UnmarshalTOML(data []byte) (err error) {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, ip := range arr {
		ip = strings.TrimSpace(ip)
		ip = strings.Trim(ip, "\"")
		if ip == "" {
			continue
		}

		ipobj := net.ParseIP(ip)
		if ipobj == nil {
			return fmt.Errorf("Bad IP format: %s", ip)
		} else {
			*self = append(*self, ipobj)
		}
	}

	self.Sort()
	return nil
}

type IPRange struct {
	Subnet *net.IPNet
	Start  uint32
	End    uint32
}
type IPNetList []IPRange

func (a IPNetList) Len() int {
	return len(a)
}
func (a IPNetList) Less(i, j int) bool {
	return a[i].End < a[j].End
}
func (a IPNetList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a IPNetList) Sort() {
	sort.Sort(a)
}

func (a IPNetList) Contains(ip net.IP) bool {
	ipval := IPToInt(ip.To4())

	l := len(a)
	i := sort.Search(l, func(i int) bool {
		n := a[i]
		return n.End >= ipval
	})

	if i < l {
		n := a[i]
		if n.Start <= ipval {
			return true
		}
	}
	return false
}

func (self *IPNetList) UnmarshalTOML(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, subnet := range arr {
		subnet = strings.TrimSpace(subnet)
		subnet = strings.Trim(subnet, "\"")
		if subnet == "" {
			continue
		}

		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			fmt.Println("ERR!")
			return err
		} else {
			start := IPToInt(ipNet.IP)
			end := start + ^IPToInt(net.IP(ipNet.Mask))

			*self = append(*self, IPRange{Subnet:ipNet, Start:start, End:end})
		}
	}

	self.Sort()
	return nil
}









