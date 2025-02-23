package main

import (
	"findservers/cache"
	"findservers/steam"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
}

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
					serverCache.MergeServers(servers)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	// Setup API routes
	r := gin.Default()

	// Serve static files
	// Serve static files
	r.Static("/static", "./static")

	// Serve index.html for root path
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	// API routes
	r.GET("/api/servers", func(c *gin.Context) {
		servers := serverCache.GetServers()
		c.JSON(200, servers)
	})

	r.Run(":8080")
}
