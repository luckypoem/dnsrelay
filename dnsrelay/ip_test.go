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

func TestIPList_FindIP(t *testing.T) {
	ips := IPList{
		net.IP{1, 1, 1, 1},
		net.IP{2, 2, 2, 2},
		net.IP{114, 114, 114, 114},
		net.IP{8, 8, 4, 4},
	}
	ips.Sort()

	for _, ip := range ips {
		if ips.FindIP(ip) != true {
			t.Fatalf("Ip %v is not in %v", ip, ips)
		}
	}
}