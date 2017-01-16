package dnsrelay

import (
	"github.com/naoina/toml"
	"os"
	"io/ioutil"
	"github.com/FTwOoO/go-logger"
	"github.com/FTwOoO/vpncore/net/rule"
	"github.com/FTwOoO/vpncore/net/addr"
	"github.com/miekg/dns"
	"net"
	"strconv"
)

var DefaultConfig *Config

type LogConfig  struct {
	LogLevel logger.LogLevel `toml:"log-level"`
	LogFile  string `toml:"log-file"`
}

type GeoIPValidate struct {
	Enable      bool     `toml:"enable"`
	Groups      []string `toml:"groups"`
	GeoIpDBPath string   `toml:"geoip-mmdb-file"`
	GeoCountry  string   `toml:"geoip-country"`
}

type Config struct {
	Addr          string   `toml:"addr"`
	DefaultGroups []string `toml:"default-group"`

	GeoIPValidate GeoIPValidate `toml:"GeoIPValidate"`

	IPFilter      rule.IPBlocker `toml:"IPFilter"`

	DNSCache      DNSCache  `toml:"Cache"`

	DNSGroups     map[string][]addr.DNSAddresss `toml:"DNSGroup"`

	DomainRules   rule.DomainRules `toml:"DomainRule"`

	Hosts         addr.Hosts `toml:"Host"`

	LogConfig     LogConfig  `toml:"Log"`
}

func NewConfig(path string) (c *Config, err error) {

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	var config Config
	if err = toml.Unmarshal(buf, &config); err != nil {
		return
	}

	return &config, nil

}

func init() {
	dnsconf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		panic(err)
	}

	systemDnsServerIp := net.ParseIP(dnsconf.Servers[0])
	systemDnsPort, _ := strconv.ParseInt(dnsconf.Port, 10, 16)

	DefaultConfig = &Config{
		DefaultGroups:[]string{rule.SYSTEM_GROUP},
		GeoIPValidate:GeoIPValidate{Enable:false},
		DNSCache:DNSCache{Backend:"memory", MinExpire:60, MaxCount:500},
		DNSGroups:map[string][]addr.DNSAddresss{rule.SYSTEM_GROUP:[]addr.DNSAddresss{
			{Ip:systemDnsServerIp,
				Port:uint16(systemDnsPort)},
		}},
		LogConfig:LogConfig{LogLevel:logger.DEBUG},
	}
}