package dnsrelay

import "net"

const CN_GROUP  = "CN"
const FG_GROUP = "FG"
const REJECT_GROUP = "REJECT"


type Config struct {
	GeoIPDBPath string

	DNSGroups     map[string] []net.IP
	DefaultGroups []string
	DomainRules    DomainRules

	BlackIP      BlackIP
}
