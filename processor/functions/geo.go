package functions

import (
	"bytes"
	"fmt"
	"net"
	"nginx-log-collector/geodb"

	"github.com/oschwald/geoip2-golang"
)

type geo struct {
	StoreTo map[string]int `yaml:"store_to"`
}

func (f *geo) Call(value string) FunctionResult {
	result := make(FunctionResult, 0, len(f.StoreTo))
	geodb := geodb.GetGeoDB()
	ip := stringToIPv4(value)
	lat, long, countryName, countryIsoCode, err := getCityData(geodb.City, ip)
	if err != nil {
		fmt.Println(err)
		return result
	}

	asnID, asnName, err := getASNData(geodb.ASN, ip)
	if err != nil {
		fmt.Println(err)
		return result
	}

	for fieldName := range f.StoreTo {
		b := bytes.Buffer{}
		dstFieldName := fieldName

		b.WriteByte('"')
		switch fieldName {
		case "lat":
			b.WriteString(fmt.Sprintf("%f", lat))
		case "long":
			b.WriteString(fmt.Sprintf("%f", long))
		case "country_name":
			b.WriteString(countryName)
		case "country_code":
			b.WriteString(countryIsoCode)
		case "asn_id":
			b.WriteString(fmt.Sprint(asnID))
		case "asn_name":
			b.WriteString(asnName)
		}
		b.WriteByte('"')

		result = append(result, FunctionPartialResult{
			Value:        b.Bytes(),
			DstFieldName: &dstFieldName,
		})
	}

	return result
}

func stringToIPv4(ip string) net.IP {
	return net.ParseIP(ip)
}

func getASNData(db *geoip2.Reader, ip net.IP) (uint, string, error) {
	record, err := db.ASN(ip)
	if err != nil {
		return 0, "", err
	}

	return record.AutonomousSystemNumber, record.AutonomousSystemOrganization, nil
}

func getCityData(db *geoip2.Reader, ip net.IP) (float64, float64, string, string, error) {
	record, err := db.City(ip)
	if err != nil {
		return 0, 0, "", "", err
	}

	return record.Location.Latitude, record.Location.Longitude, record.Country.Names["en"], record.Country.IsoCode, nil
}
