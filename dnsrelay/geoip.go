//
// Modified from https://github.com/oschwald/geoip2-golang/reader.go
//
// why not use geoip2-golang ?
// 	If you only need several fields, you may get superior performance by using maxminddb's
// 	Lookup directly with a result struct that only contains the required fields.
// 	(See example_test.go in the https://github.com/oschwald/maxminddb-golang repository for an example of this.)
//

package dnsrelay

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

// The Country struct corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
//
// See http://dev.maxmind.com/geoip/geoip2/whats-new-in-geoip2/ for example data and description:
// The country is the country where the IP address is located.
// The registered_country is the country in which the IP is registered. These two may differ in some cases.
//


type Country struct {
	Country struct {
		IsoCode   string            `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		IsoCode   string            `maxminddb:"iso_code"`
	} `maxminddb:"registered_country"`
}



type databaseType int

const (
	isAnonymousIP = 1 << iota
	isCity
	isConnectionType
	isCountry
	isDomain
	isEnterprise
	isISP
)

// Reader holds the maxminddb.Reader struct. It can be created using the
// Open and FromBytes functions.
type Reader struct {
	mmdbReader   *maxminddb.Reader
	databaseType databaseType
}

// InvalidMethodError is returned when a lookup method is called on a
// database that it does not support. For instance, calling the ISP method
// on a City database.
type InvalidMethodError struct {
	Method       string
	DatabaseType string
}

func (e InvalidMethodError) Error() string {
	return fmt.Sprintf(`geoip2: the %s method does not support the %s database`,
		e.Method, e.DatabaseType)
}

// UnknownDatabaseTypeError is returned when an unknown database type is
// opened.
type UnknownDatabaseTypeError struct {
	DatabaseType string
}

func (e UnknownDatabaseTypeError) Error() string {
	return fmt.Sprintf(`geoip2: reader does not support the "%s" database type`,
		e.DatabaseType)
}

// Open takes a string path to a file and returns a Reader struct or an error.
// The database file is opened using a memory map. Use the Close method on the
// Reader object to return the resources to the system.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}


func getDBType(reader *maxminddb.Reader) (databaseType, error) {
	switch reader.Metadata.DatabaseType {
	case "GeoIP2-Anonymous-IP":
		return isAnonymousIP, nil
	// We allow City lookups on Country for back compat
	case "GeoLite2-City", "GeoIP2-City", "GeoIP2-Precision-City", "GeoLite2-Country",
		"GeoIP2-Country":
		return isCity | isCountry, nil
	case "GeoIP2-Connection-Type":
		return isConnectionType, nil
	case "GeoIP2-Domain":
		return isDomain, nil
	case "GeoIP2-Enterprise":
		return isEnterprise | isCity | isCountry, nil
	case "GeoIP2-ISP", "GeoIP2-Precision-ISP":
		return isISP, nil
	default:
		return 0, UnknownDatabaseTypeError{reader.Metadata.DatabaseType}
	}
}

// Country takes an IP address as a net.IP struct and returns a Country struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 Country
// databases.
func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	if isCountry&r.databaseType == 0 {
		return nil, InvalidMethodError{"Country", r.Metadata().DatabaseType}
	}
	var country Country
	err := r.mmdbReader.Lookup(ipAddress, &country)
	return &country, err
}

// Metadata takes no arguments and returns a struct containing metadata about
// the MaxMind database in use by the Reader.
func (r *Reader) Metadata() maxminddb.Metadata {
	return r.mmdbReader.Metadata
}

// Close unmaps the database file from virtual memory and returns the
// resources to the system.
func (r *Reader) Close() error {
	return r.mmdbReader.Close()
}


func (r *Reader) IsChineseIP(ipAddress net.IP) (bool, error) {
	country, err := r.Country(ipAddress)
	if err != nil {
		return false, err
	}

	return country.Country.IsoCode == "CN" || country.RegisteredCountry.IsoCode == "CN", nil
}