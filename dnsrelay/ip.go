package dnsrelay

import (
	"net"
	"encoding/binary"
	"strings"
	"fmt"
	"strconv"
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

func (self *IPList) UnmarshalTOML(data []byte) (err error) {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, ip := range arr {
		ip = strings.TrimSpace(ip)
		ip = strings.Trim(ip, "\"")

		ipobj := net.ParseIP(ip)
		if ipobj == nil {
			return fmt.Errorf("Bad IP format: %s", ip)
		} else {
			*self = append(*self, ipobj)
		}
	}
	return nil
}

type IPNetList []net.IPNet

func (self *IPNetList) UnmarshalTOML(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, subnet := range arr {
		subnet = strings.TrimSpace(subnet)
		subnet = strings.Trim(subnet, "\"")

		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return err
		} else {
			*self = append(*self, *ipNet)
		}
	}
	return nil
}


