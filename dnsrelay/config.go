package dnsrelay

import (
	"github.com/naoina/toml"
	"os"
	"io/ioutil"
	"net"
	"strings"
	"errors"
	"strconv"
	"fmt"
)

const CN_GROUP = "CN"
const FG_GROUP = "FG"
const REJECT_GROUP = "REJECT"

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
			return errors.New("Bad format for DNS")
		}
		ip, port := arr[0], arr[1]
		self.Ip = net.ParseIP(ip)
		self.Port, err = strconv.Atoi(port)
		if err != nil {
			return err
		}
	} else {
		self.Ip = net.ParseIP(s)
		self.Port = 53
	}
	return nil
}

func (self *DNSAddresss) String() string {
	return fmt.Sprintf("%s:%d", self.Ip.String(), self.Port)
}

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

	DNSGroups     map[string][]DNSAddresss `toml:"DNSGroup"`

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