/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package dnsrelay

import (
	"testing"
	"net"
)

func TestIPList_Contains(t *testing.T) {
	ips := IPList{
		net.IP{1, 1, 1, 1},
		net.IP{2, 2, 2, 2},
		net.IP{114, 114, 114, 114},
		net.IP{8, 8, 4, 4},
	}
	ips.Sort()

	for _, ip := range ips {
		if ips.Contains(ip) != true {
			t.Fatalf("Ip %v is not in %v", ip, ips)
		}
	}
}

func TestIPNetList_Contains(t *testing.T) {
	ips := IPList{
		net.IP{32, 32, 32, 5},
		net.IP{17, 3, 2, 2},
		net.IP{22, 33, 44, 254},
		net.IP{22, 33, 44, 253},
	}

	netl := new(IPNetList)
	err := netl.UnmarshalTOML([]byte(`[
	    "32.32.32.0/24", "17.3.4.2/16", "22.33.44.253/30",
	    ]`))

	if err != nil {
		t.Fatal(err)
	}

	for _, ip := range ips {
		if netl.Contains(ip) != true {
			t.Fatalf("Ip %v is not in %v", ip, netl)
		}
	}

}