package cache

import (
	"findservers/models"
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

func (c *ServerCache) MergeServers(newServers []models.Server) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create map of existing servers for O(1) lookup
	existingServers := make(map[string]models.Server)
	for _, server := range c.servers {
		existingServers[server.Addr] = server
	}

	now := time.Now()

	// Merge/update servers
	for _, newServer := range newServers {
		if existing, exists := existingServers[newServer.Addr]; exists {
			// Update existing server
			newServer.FirstSeen = existing.FirstSeen
			newServer.LastSeen = now
			existingServers[newServer.Addr] = newServer
		} else {
			// Add new server
			newServer.FirstSeen = now
			newServer.LastSeen = now
			existingServers[newServer.Addr] = newServer
		}
	}

	// Update last seen time for all existing servers that were returned
	for addr, server := range existingServers {
		if _, found := existingServers[addr]; found {
			server.LastSeen = now
			existingServers[addr] = server
		}
	}

	// Convert map back to slice
	c.servers = make([]models.Server, 0, len(existingServers))
	for _, server := range existingServers {
		c.servers = append(c.servers, server)
	}

	c.lastUpdate = now
}

func (c *ServerCache) GetServers() []models.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.servers
}

func (c *ServerCache) NeedsUpdate() bool {
	return time.Since(c.lastUpdate) > 5*time.Minute
}

func (c *ServerCache) PruneInactiveServers(maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	activeServers := make([]models.Server, 0)

	for _, server := range c.servers {
		if now.Sub(server.LastSeen) < maxAge {
			activeServers = append(activeServers, server)
		}
	}

	c.servers = activeServers
}
