package traefik_geoipstate

import (
	"context"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

type Config struct {
	Database string `json:"database"`
}

func CreateConfig() *Config {
	return &Config{}
}

type GeoIPState struct {
	next http.Handler
	db   *geoip2.Reader
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	db, err := geoip2.Open(config.Database)
	if err != nil {
		return nil, err
	}

	return &GeoIPState{
		next: next,
		db:   db,
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

	record, err := m.db.City(ip)
	if err != nil {
		m.next.ServeHTTP(rw, req)
		return
	}

	// Country
	if record.Country.IsoCode != "" {
		req.Header.Set("X-User-Country", record.Country.IsoCode)
	}

	// State (US subdivision)
	if len(record.Subdivisions) > 0 && record.Subdivisions[0].IsoCode != "" {
		req.Header.Set("X-User-State", record.Subdivisions[0].IsoCode)
	}

	m.next.ServeHTTP(rw, req)
}
