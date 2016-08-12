# dnsrelay

dnsrelay is a DNS proxy like [godns](https://github.com/kenshinx/godns)and [ChinaDNS](https://github.com/shadowsocks/ChinaDNS). The goal of this project is to escape from DNS poisoning powered by GFW(The Great Firewall   Of China)

Thans to [godns](https://github.com/kenshinx/godns),[grimd](https://github.com/looterz/grimd),[ChinaDNS](https://github.com/shadowsocks/ChinaDNS),[dnsserver](https://github.com/docker/dnsserver),[dns-reverse-proxy](https://github.com/StalkR/dns-reverse-proxy) for the idea.

# depandencies

* [GeoLite2 Free Downloadable Databases](http://dev.maxmind.com/geoip/geoip2/geolite2/)
* [geoip2-golang](https://github.com/oschwald/geoip2-golang) and [maxminddb-golang](https://github.com/oschwald/maxminddb-golang), open source Third-Party [MaxMind](http://maxmind.github.io/MaxMind-DB/) Reader for Golang
* [dns](https://github.com/miekg/dns)
* [toml](https://github.com/naoina/toml), [TOML config file](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.4.0.md) parsing 
* [go-logger](https://github.com/apsdehal/go-logger)

## Feature/TODO
1. Query multiple upstream DNS servers group concurrently
2. Cache all mostly used domain names
3. Load all mostly used domain names at startup
4. Map hostname to ip statically by configuration
5. Domain name matching for custom DNS server
6. GeoIP strategy for filtering untrusted DNS results from DNS server in China 
7. Black IP list for filtering untrusted DNS results

## Configuration

The configuration dnsrelay.toml is a [TOML](https://github.com/mojombo/toml) format config file.

## Notice
If DNS protocol are poisoning and filtering like in  China, DNS server like 8.8.8.8 may not response, so VPN(and system routing tables entry for 8.8.8.8, e.g.) is required to get dnsrelay work.


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
