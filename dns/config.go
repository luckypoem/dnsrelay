package dns

import (
	"github.com/naoina/toml"
	"github.com/apsdehal/go-logger"
	"github.com/FTwOoO/vpncore/net/rule"
	"github.com/FTwOoO/vpncore/net/addr"
	"github.com/miekg/dns"
	"net"
	"strconv"
	"errors"
	"os"
	"io/ioutil"
)

var (
	ErrDnsTimeOut = errors.New("dns timeout.")
	ErrDnsMsgIllegal = errors.New("dns message illegal.")
	ErrNoDnsServer = errors.New("no proper dns server.")
)

var DefaultConfig *Config

func init() {
	dnsconf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	var dnsAddrs = []addr.DNSAddresss{}

	if err == nil {
		systemDnsServerIp := net.ParseIP(dnsconf.Servers[0])
		systemDnsPort, _ := strconv.ParseInt(dnsconf.Port, 10, 16)
		addr1 := addr.DNSAddresss{Ip:systemDnsServerIp, Port:uint16(systemDnsPort)}
		dnsAddrs = append(dnsAddrs, addr1)
	} else {
		dnsAddrs = append(dnsAddrs, addr.DNSAddresss{Ip:net.IP{8, 8, 8, 8}, Port:uint16(53)})
		dnsAddrs = append(dnsAddrs, addr.DNSAddresss{Ip:net.IP{114, 114, 114, 114}, Port:uint16(53)})
	}

	DefaultConfig = &Config{
		Addr:":5353",
		DefaultGroups:[]string{rule.SYSTEM_GROUP},
		GeoIPValidate:GeoIPValidateConfig{Enable:false},
		DNSCache:DNSCacheConfig{Enable:true, MaxCount:500},
		DNSGroups:map[string][]addr.DNSAddresss{rule.SYSTEM_GROUP:dnsAddrs},
		LogConfig:LogConfig{LogLevel:logger.DEBUG},
	}
}

type LogConfig  struct {
	LogLevel logger.LogLevel `toml:"log-level"`
	LogFile  string `toml:"log-file"`
}

type GeoIPValidateConfig struct {
	Enable      bool     `toml:"enable"`
	Groups      []string `toml:"groups"`
	GeoIpDBPath string   `toml:"geoip-mmdb-file"`
	GeoCountry  string   `toml:"geoip-country"`
}

type Config struct {
	Addr          string   `toml:"addr"`
	DefaultGroups []string `toml:"default-group"`

	GeoIPValidate GeoIPValidateConfig `toml:"GeoIPValidate"`

	IPBlocker     rule.IPBlocker `toml:"IPBlocker"`

	DNSCache      DNSCacheConfig  `toml:"Cache"`

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

