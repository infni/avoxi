package countycodes

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type CountryCodeComparer interface {
	IsInCountry(ip string, countries []string) (bool, error)
}

func NewContryCodeComparer(dbFileLocation string) CountryCodeComparer {
	comparer := new(countryCodeComparerImpl)
	comparer.dbFileLocation = dbFileLocation
	return comparer
}

type countryCodeComparerImpl struct {
	dbFileLocation string
}

func (comparer *countryCodeComparerImpl) IsInCountry(ip string, countries []string) (bool, error) {
	db, dbErr := maxminddb.Open(comparer.dbFileLocation) // "test-data/test-data/GeoIP2-City-Test.mmdb"
	if dbErr != nil {
		return false, fmt.Errorf("Failed to open local mmdb : %w", dbErr)
	}
	defer db.Close()

	var record struct {
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	} // Or any appropriate struct

	if netip := net.ParseIP(ip); netip == nil {
		// if the IP is malformed return false AND an eror.
		return false, fmt.Errorf("IP is malformed : %s", ip)
	} else if err := db.Lookup(netip, &record); err != nil {
		// if the record is missing (IP is not in the db), default to returning false AND an eror.
		return false, fmt.Errorf("Failed to lookup record : %w", err)
	}

	for _, inCountry := range countries {
		for _, compCountry := range record.Country.Names {
			if inCountry == compCountry {
				return true, nil
			}
		}
	}

	return false, nil
}
