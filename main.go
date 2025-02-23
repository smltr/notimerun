package main

import (
	"findservers/cache"
	"findservers/models"
	"findservers/steam"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
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
	r.Use(gzip.Gzip(gzip.DefaultCompression))

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

	r.GET("/api/servers/list", func(c *gin.Context) {
		servers := serverCache.GetServers()
		c.Header("Content-Type", "text/html")
		c.String(200, generateServerHTML(servers))
	})

	r.Run(":8080")
}

func generateServerHTML(servers []models.Server) string {
	var html strings.Builder

	html.WriteString(fmt.Sprintf(`<span id="server-count" hx-swap-oob="true">%d</span>`, len(servers)))

	for _, server := range servers {
		// Create a unique ID using the server address (you might want to hash this if addresses are too long)
		serverID := fmt.Sprintf("server-%s", strings.Replace(server.Addr, ":", "-", -1))

		secure := ""
		if server.Secure {
			secure = "●"
		}
		bots := ""
		if server.Bots > 0 {
			bots = fmt.Sprintf("%d", server.Bots)
		}
		html.WriteString(fmt.Sprintf(`
            <div id="%s"
                class="server-row"
                x-data
                :class="{ 'selected': selectedServer === '%s' }"
                @click="selectedServer = '%s'"
                @dblclick="window.location.href = 'steam://connect/%s'">
                <div class="cell-pw"></div>
                <div class="cell-vac">%v</div>
                <div class="cell-region">%s</div>
                <div class="cell-name">
                    <span class="server-name">%s</span>
                    <span class="server-ip">%s</span>
                </div>
                <div class="cell-bot">%s</div>
                <div class="cell-players">%s</div>
                <div class="cell-map">%s</div>
                <div class="cell-tags">%s</div>
            </div>`,
			serverID, // Add the unique ID here
			server.Addr,
			server.Addr,
			server.Addr,
			secure,
			regionCodeToString(server.Region),
			template.HTMLEscapeString(server.Name),
			server.Addr,
			bots,
			formatPlayers(server.Players, server.MaxPlayers),
			template.HTMLEscapeString(server.Map),
			template.HTMLEscapeString(server.GameType),
		))
	}
	return html.String()
}

func regionCodeToString(code int) string {
	regions := map[int]string{
		0: "US", 1: "US", 2: "SA", 3: "EU",
		4: "AS", 5: "AU", 6: "ME", 7: "AF",
		255: "WD",
	}
	if reg, ok := regions[code]; ok {
		return reg
	}
	return "??"
}

func formatPlayers(players, maxPlayers int) string {
	if players == 0 {
		return fmt.Sprintf("<span class='zero-players'>%d</span><span class='max-players'>/%d</span>",
			players, maxPlayers)
	}
	return fmt.Sprintf("%d<span class='max-players'>/%d</span>", players, maxPlayers)
}
