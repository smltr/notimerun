package main

import (
	"fmt"
	"log"
	"notimerun/cache"
	"notimerun/steam"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	r.StaticFile("/service-worker.js", "./static/service-worker.js")

	// API routes
	r.GET("/api/servers", func(c *gin.Context) {
		servers := serverCache.GetServers()
		c.JSON(200, servers)
	})

	r.POST("/api/report", func(c *gin.Context) {
		var report ReportRequest
		if err := c.BindJSON(&report); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Create email
		from := mail.NewEmail("CS2 Server Browser", "noreply@notime.run")
		to := mail.NewEmail("Sam", "smltr@proton.me") // Replace with your email
		subject := fmt.Sprintf("New %s Report", report.Type)
		content := fmt.Sprintf(
			"Type: %s\nDescription: %s\nContact: %s",
			report.Type,
			report.Description,
			report.Contact,
		)
		message := mail.NewSingleEmail(from, subject, to, content, content)

		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		_, err := client.Send(message)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to send report"})
			return
		}

		c.JSON(200, gin.H{"message": "Report submitted successfully"})
	})

	r.Run(":8080")
}
