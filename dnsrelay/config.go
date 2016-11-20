package dnsrelay

import (
	"github.com/naoina/toml"
	"os"
	"io/ioutil"
	"github.com/FTwOoO/go-logger"
	"github.com/FTwOoO/vpncore/net/rule"
	"github.com/FTwOoO/vpncore/net/addr"
)

type Config struct {
	GeoIPDBPath   string   `toml:"geoip-mmdb-file"`
	ADDR          string   `toml:"addr"`
	PORT          uint     `toml:"port"`
	CacheNum      uint     `toml:"cache-num"`
	FuckGFW       bool     `toml:"fuck-gfw"`
	DefaultGroups []string `toml:"default-group"`
	LogFile       string   `toml:"log-file"`

	IPFilter      rule.IPBlocker `toml:"IPFilter"`

	DNSCache      DNSCache  `toml:"Cache"`

	DNSGroups     map[string][]addr.DNSAddresss `toml:"DNSGroup"`

	DomainRules   rule.DomainRules `toml:"DomainRule"`

	Hosts         addr.Hosts `toml:"Host"`

	// these fields are not fields from config file
	LogConfig     struct {
			      LogLevel logger.LogLevel `toml:"log-level"`
			      LogFile  string `toml:"log-file"`
		      } `toml:"Log"`
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