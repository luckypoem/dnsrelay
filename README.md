# dnsrelay

dnsrelay is a DNS proxy like [godns](https://github.com/kenshinx/godns) and [ChinaDNS](https://github.com/shadowsocks/ChinaDNS). The goal of this project is to escape from DNS poisoning powered by GFW(The Great Firewall   Of China)

Thans to [godns](https://github.com/kenshinx/godns),[grimd](https://github.com/looterz/grimd),[ChinaDNS](https://github.com/shadowsocks/ChinaDNS),[dnsserver](https://github.com/docker/dnsserver),[dns-reverse-proxy](https://github.com/StalkR/dns-reverse-proxy) for the idea.

# depandencies

* [GeoLite2 Free Downloadable Databases](http://dev.maxmind.com/geoip/geoip2/geolite2/)
* [geoip2-golang](https://github.com/oschwald/geoip2-golang) and [maxminddb-golang](https://github.com/oschwald/maxminddb-golang), open source Third-Party [MaxMind](http://maxmind.github.io/MaxMind-DB/) Reader for Golang
* [go dns](https://github.com/miekg/dns)
* [toml](https://github.com/naoina/toml), [TOML config file](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.4.0.md) parsing 
* [go-logger](https://github.com/apsdehal/go-logger)

## Feature
1. Query multiple upstream DNS group concurrently
2. Cache all mostly used domain names
3. Hosts configuration
4. Domain name matching for custom DNS server
5. GeoIP strategy for filtering untrusted DNS results from DNS server of China 
6. Black IP list for filtering untrusted DNS results

## TODO
1. Load all mostly used domain names at startup
2. Supporting Ip filtering by subnet config

## Notice
If DNS protocol are poisoning and filtering like in  China, DNS server like 8.8.8.8 may not response, so VPN(and some routing tables entry for 8.8.8.8, e.g.) is required to get dnsrelay work.

## Configuration

The configuration dnsrelay.toml is a [TOML](https://github.com/mojombo/toml) format config file.

* DNS group

You can define multiple groups, default groups is used to send all DNS request that no rules match.

```
default-group = ["CN","FG"]

[DNSGroup]
CN=["192.168.1.1", "114.114.114.114", "223.6.6.6:53", "223.5.5.5"]
FG=["8.8.8.8", "8.8.4.4"]
```

* DNS matching Rules

For example, if there is a DNS request asking A record for “baidu.com”, the rule bellow matched, the DNS group specified by ‘domain-group’ is CN. Then the the DNS request is send to all DNS defined by CN group

```
[[DomainRule]]
scheme="DOMAIN-MATCH"
group="CN"
value=[
    "order.mi.com",
    "baidu.com",
    ]
```

 * Hosts

Like /etc/hosts, you can set one or more ip for a domain without sending any DNS requests:

```
[[Host]]
name="example.com"
ip=["93.184.216.34"]

[[Host]]
name="ip.com"
ip=["64.111.96.203"]
```

* Ip filtering

DNS result may points to the wrong Ip in China because of DNS poisoning power by GFW(?), so ip filtering is needed:

```
[IPFilter]
ip=[
   "64.111.96.204"
   ]
```

* Cache DNS results

```
[Cache]
backend = "memory"
expire = 3600 
maxcount = 500
```


## LICENSE

```
Copyright (c) 2016 <booopooob@gmail.com>

This program is free software: you can redistribute it and/or modify    
it under the terms of the GNU General Public License as published by    
the Free Software Foundation, either version 3 of the License, or    
(at your option) any later version.    

This program is distributed in the hope that it will be useful,    
but WITHOUT ANY WARRANTY; without even the implied warranty of    
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the    
GNU General Public License for more details.    

You should have received a copy of the GNU General Public License    
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
