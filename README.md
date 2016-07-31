# DNSPLAY

DNSPLAY is a DNS proxy like [dns-reverse-proxy](https://github.com/StalkR/dns-reverse-proxy) and [ShadowDNS](https://github.com/shadowsocks/ShadowDNS). The goal of this project is to escape from DNS poisoning powered by GFW(The Great Firewall   Of China)

Thans to  [dns-reverse-proxy](https://github.com/StalkR/dns-reverse-proxy) 、[ShadowDNS](https://github.com/shadowsocks/ShadowDNS)、[jianbing-dictionary-dns](https://github.com/chuangbo/jianbing-dictionary-dns/blob/master/golang/jianbing-dns/jianbing-dns.go)、[dnsserver](https://github.com/docker/dnsserver) and the powerful go [dns](https://github.com/miekg/dns) library.

## Feature/TODO
1. Query multiple upstream DNS servers concurrently
2. Cache all mostly used domain names
3. Load all mostly used domain names at startup
4. DNS server group
5. Domain name matching for custom DNS server
6. GeoIP strategy for filtering untrusted DNS results from DNS server in China 
7. Black IP list for filtering untrusted DNS results

## Notice
If DNS protocol are poisoning and filtering like in  China, DNS server like 8.8.8.8 may not response, so VPN(and system routing tables entry for 8.8.8.8, e.g.) is required to get DNSPLAY work.




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
