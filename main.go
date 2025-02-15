package main

import (
	"frfr/cache"
	"frfr/steam"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	steamClient := steam.NewSteamClient()
	serverCache := cache.NewServerCache()

	// Start background server list updater
	go func() {
		for {
			if serverCache.NeedsUpdate() {
				servers, err := steamClient.FetchServers()
				if err != nil {
					log.Printf("Error fetching servers: %v", err)
				} else {
					serverCache.UpdateServers(servers)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	// Setup API routes
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API routes
	r.GET("/api/servers", func(c *gin.Context) {
		servers := serverCache.GetServers()
		c.JSON(200, servers)
	})

	r.Run(":8080")
}
