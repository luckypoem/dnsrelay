package dnsrelay

import (
	"testing"
	"time"
	"fmt"
	"net"
	"math/rand"
)

func BenchmarkMaxMindDBCountry(b *testing.B) {
	start := time.Now()
	fmt.Printf("\nTest %d times:\n", b.N)

	db, err := Open("GeoLite2-Country.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	r := rand.New(rand.NewSource(0))

	for i := 0; i < b.N; i++ {
		ip := randomIPv4Address(b, r)
		isCN, err := db.IsChineseIP(ip)
		country, _ := db.Country(ip)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Printf("IP %v is CN IP ? : %v [%s:%s]\n", ip, isCN,
			country.Country.IsoCode, country.RegisteredCountry.IsoCode)
	}
	fmt.Printf("Time used: %s\n", time.Since(start))

}

func randomIPv4Address(b *testing.B, r *rand.Rand) net.IP {
	num := r.Uint32()
	return []byte{byte(num >> 24), byte(num >> 16), byte(num >> 8),
		byte(num)}
}
