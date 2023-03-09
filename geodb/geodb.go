package geodb

import (
	"sync"

	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog"
)

type GeoDB struct {
	ASN, City *geoip2.Reader
}

func GetGeoDB() *GeoDB {
	return geoDBInstance
}

func InitGeoDB(city, asn string, logger *zerolog.Logger) {
	once.Do(func() {
		geoDBInstance = &GeoDB{}
		var err error
		geoDBInstance.City, err = geoip2.Open(city)
		if err != nil {
			logger.Fatal().Err(err).Msg("unable to load city database")
		}

		geoDBInstance.ASN, err = geoip2.Open(asn)
		if err != nil {
			logger.Fatal().Err(err).Msg("unable to load asn database")
		}
	})
}

var geoDBInstance *GeoDB
var once sync.Once
