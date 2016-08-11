package dnsrelay

import (
	"github.com/naoina/toml"
	"os"
	"io/ioutil"
	"net"
)

const CN_GROUP = "CN"
const FG_GROUP = "FG"
const REJECT_GROUP = "REJECT"


type IPFilter struct {
	Ip  [] string
	Net [] string
}

type Config struct {
	GeoIPDBPath   string `toml:"geoip-mmdb-file"`
	CacheNum      uint        `toml:"cache-num"`
	IsInChina     bool        `toml:"in-china"`
	DefaultGroups []string `toml:"default-group"`
	LogFile       string `toml:"log-file"`

	IPFilter      struct {
			      Ip  [] string
			      Net [] string
		      }         `toml:"IPFilter"`

	DNSGroups     map[string][]string `toml:"DNSGroup"`

	Rules         []DomainRule `toml:"DomainRule"`

	Hosts         [] struct {
		Name string
		Ip   string
	}       `toml:"Host"`

	// these fields are not fields from config file
	IPBlocker     *IPBlocker
	Logger        *Logger
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

	ips := IPList{}

	for _, ip := range config.IPFilter.Ip {
		ips = append(ips, net.ParseIP(ip))
	}

	config.IPBlocker = NewIPBlocker(ips)

	config.Logger, err = NewLogger(config.LogFile, "dnsrelay")
	if err != nil {
		return
	}

	return &config, nil

}