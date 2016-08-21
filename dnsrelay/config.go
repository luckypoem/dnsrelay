package dnsrelay

import (
	"github.com/naoina/toml"
	"os"
	"io/ioutil"
	"github.com/FTwOoO/go-logger"

)

type Config struct {
	GeoIPDBPath   string   `toml:"geoip-mmdb-file"`
	CacheNum      uint     `toml:"cache-num"`
	FuckGFW       bool     `toml:"fuck-gfw"`
	DefaultGroups []string `toml:"default-group"`
	LogFile       string   `toml:"log-file"`

	IPFilter      IPFilter `toml:"IPFilter"`

	DNSCache      DNSCache  `toml:"Cache"`

	DNSGroups     map[string][]DNSAddresss `toml:"DNSGroup"`

	DomainRules   DomainRules `toml:"DomainRule"`

	Hosts         Hosts `toml:"Host"`

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