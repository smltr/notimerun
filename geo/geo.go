package geo

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

type Locator struct {
	db    *geoip2.Reader
	mu    sync.RWMutex
	cache map[string]*geoip2.City
}

func NewLocator(dbPath string) (*Locator, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &Locator{
		db:    db,
		cache: make(map[string]*geoip2.City),
	}, nil
}

func (l *Locator) GetLocation(ipStr string) (*geoip2.City, error) {
	// Strip port if present
	host, _, err := net.SplitHostPort(ipStr)
	if err != nil {
		host = ipStr
	}

	// Check cache first
	l.mu.RLock()
	if location, ok := l.cache[host]; ok {
		l.mu.RUnlock()
		return location, nil
	}
	l.mu.RUnlock()

	// Parse IP and look up
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", host)
	}

	record, err := l.db.City(ip)
	if err != nil {
		return nil, err
	}

	// Cache the result
	l.mu.Lock()
	l.cache[host] = record
	l.mu.Unlock()

	return record, nil
}

func (l *Locator) Close() error {
	return l.db.Close()
}
