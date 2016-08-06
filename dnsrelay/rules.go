package dnsrelay

import (
	"strings"
	"net"
	"sort"
)

const (
	MATCH_TYPE_DOMAIN_SUFFIX = iota
	MATCH_TYPE_DOMAIN
	MATCH_TYPE_DOMAIN_KEYWORD
)

type Rule struct {
	Group     string
	MatchType uint
	Value     string
}

func (rule Rule) Match(input string) bool {
	switch rule.MatchType {
	case MATCH_TYPE_DOMAIN:
		return input == rule.Value
	case MATCH_TYPE_DOMAIN_SUFFIX:
		return strings.HasSuffix(input, rule.Value)
	case MATCH_TYPE_DOMAIN_KEYWORD:
		return strings.Contains(input, rule.Value)
	default:
		return false

	}
}

type DomainRules struct {
	Rules []Rule
}

func (self *DomainRules) AddRule(rule Rule) {
	self.Rules = append(self.Rules, rule)
}

func (self *DomainRules) FindGroup(domain string) string {
	for _, rule := range self.Rules {
		if rule.Match(domain) {
			return rule.Group
		}
	}
	return ""
}


func IPToU32(ip net.IP) uint32 {
	return (uint32)(ip[0]) << 24 | (uint32)(ip[1]) << 16 | (uint32)(ip[2]) << 8 | (uint32)(ip[3])
}


type IPList []net.IP
func (a IPList) Len() int           { return len(a) }
func (a IPList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IPList) Less(i, j int) bool {
	return IPToU32(a[i]) < IPToU32(a[j])
}

type BlackIP struct{
	ips IPList
}

func NewBlackIP (ips IPList) (*BlackIP) {
	sort.Sort(ips)
	return &BlackIP{ips:ips}
}

func (self *BlackIP) FindIP(ip net.IP) bool {
    	i := sort.Search(len(self.ips), func(i int) bool { return IPToU32(self.ips[i]) >= IPToU32(ip)})
   	return i < len(self.ips) && self.ips[i].Equal(ip)
}