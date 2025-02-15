package cache

import (
	"notimerun/models"
	"sync"
	"time"
)

type ServerCache struct {
	servers    []models.Server
	mu         sync.RWMutex
	lastUpdate time.Time
}

func NewServerCache() *ServerCache {
	return &ServerCache{}
}

func (c *ServerCache) UpdateServers(servers []models.Server) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.servers = servers
	c.lastUpdate = time.Now()
}

func (c *ServerCache) GetServers() []models.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.servers
}

func (c *ServerCache) NeedsUpdate() bool {
	return time.Since(c.lastUpdate) > 5*time.Minute
}
