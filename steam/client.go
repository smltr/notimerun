package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"findservers/geo"
	"findservers/models"
)

type SteamClient struct {
	apiKey string
}

func NewSteamClient() *SteamClient {
	return &SteamClient{
		apiKey: os.Getenv("STEAM_API_KEY"),
	}
}

func (c *SteamClient) FetchServers() ([]models.Server, error) {
	var servers []models.Server
	maxRetries := 3

	locator, err := geo.NewLocator("GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("Warning: Geo location disabled: %v", err)
	}
	defer locator.Close()

	for i := 0; i < maxRetries; i++ {
		// Try both CS2 and CSGO filters
		url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&filter=\\appid\\730\\dedicated\\1\\&limit=10000", c.apiKey)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Attempt %d: HTTP error: %v", i+1, err)
			continue
		}

		log.Printf("Attempt %d: Status Code: %d", i+1, resp.StatusCode)

		rawBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			log.Printf("Attempt %d: Body read error: %v", i+1, err)
			continue
		}

		var result struct {
			Response struct {
				Servers []models.Server `json:"servers"`
			} `json:"response"`
		}

		if err := json.Unmarshal(rawBody, &result); err != nil {
			log.Printf("Attempt %d: JSON parse error: %v", i+1, err)
			log.Printf("Response preview: %s", string(rawBody[:min(1000, len(rawBody))]))
			continue
		}

		// Debug log before filtering
		log.Printf("Before filtering - Total servers: %d", len(result.Response.Servers))

		// Log some sample server data
		for i, server := range result.Response.Servers[:min(5, len(result.Response.Servers))] {
			log.Printf("Sample server %d: Name: %s, Product: %s, GameDir: %s, AppID: %d",
				i, server.Name, server.Product, server.GameDir, server.AppID)
		}

		filteredServers := []models.Server{}
		for _, server := range result.Response.Servers {
			// Skip Valve official servers
			if strings.HasPrefix(server.Name, "Valve Counter-Strike") {
				continue
			}

			// Debug log server details
			log.Printf("Processing server: %s, Addr: %s, Product: %s",
				server.Name, server.Addr, server.Product)

			if locator != nil {
				if loc, err := locator.GetLocation(server.Addr); err == nil {
					server.CountryCode = loc.Country.IsoCode
					server.CountryName = loc.Country.Names["en"]
					server.ContinentCode = loc.Continent.Code
					log.Printf("Location found for %s: Country: %s", server.Addr, server.CountryCode)
				} else {
					log.Printf("Failed to get location for %s: %v", server.Addr, err)
				}
			}

			filteredServers = append(filteredServers, server)
		}

		log.Printf("After filtering - Servers remaining: %d", len(filteredServers))

		if len(filteredServers) > 10 {
			return filteredServers, nil
		}

		time.Sleep(time.Second * 2)
	}

	return servers, fmt.Errorf("failed to fetch sufficient servers after %d attempts", maxRetries)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
