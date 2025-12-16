package traefik_geoipstate

import (
	"context"
	"net"
	"net/http"

	"github.com/IncSW/geoip2"
)

type Config struct {
	Database string `json:"database"`
}

func CreateConfig() *Config {
	return &Config{}
}

type GeoIPState struct {
	next   http.Handler
	reader *geoip2.CityReader
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	reader, err := geoip2.NewCityReaderFromFile(config.Database)
	if err != nil {
		return nil, err
	}

	return &GeoIPState{
		next:   next,
		reader: reader,
	}, nil
}

func (m *GeoIPState) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Default values
	req.Header.Set("X-User-Country", "UNKNOWN")
	req.Header.Set("X-User-State", "UNKNOWN")

	// Get client IP
	ipStr := req.Header.Get("X-Real-IP")
	if ipStr == "" {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err == nil {
			ipStr = host
		}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		m.next.ServeHTTP(rw, req)
		return
	}

	record := m.reader.Lookup(ip)
	if record == nil {
		m.next.ServeHTTP(rw, req)
		return
	}

	// Country
	if record.Country.ISOCode != "" {
		req.Header.Set("X-User-Country", record.Country.ISOCode)
	}

	// State (subdivision)
	if len(record.Subdivisions) > 0 && record.Subdivisions[0].ISOCode != "" {
		req.Header.Set("X-User-State", record.Subdivisions[0].ISOCode)
	}

	m.next.ServeHTTP(rw, req)
}
