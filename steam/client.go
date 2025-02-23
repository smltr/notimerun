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

// func logDuplicateServers(servers []models.Server) {
// 	log.Printf("Starting duplicate analysis for %d servers", len(servers))

// 	// Create a map of IP to servers
// 	ipMap := make(map[string][]models.Server)

// 	// Group servers by IP (without port)
// 	for i, server := range servers {
// 		if i%1000 == 0 { // Log progress every 1000 servers
// 			log.Printf("Processing server %d of %d", i, len(servers))
// 		}

// 		ip := strings.Split(server.Addr, ":")[0]
// 		ipMap[ip] = append(ipMap[ip], server)
// 	}

// 	log.Printf("Found %d unique IPs", len(ipMap))

// 	// Set a reasonable limit for logging
// 	duplicateCount := 0
// 	const maxDuplicatesToLog = 100

// 	// Log IPs with multiple servers
// 	for ip, serverGroup := range ipMap {
// 		if len(serverGroup) > 1 {
// 			duplicateCount++
// 			if duplicateCount > maxDuplicatesToLog {
// 				log.Printf("Reached maximum duplicate logging limit of %d. Stopping.", maxDuplicatesToLog)
// 				break
// 			}

// 			truncateString := func(s string, maxLen int) string {
// 				if len(s) <= maxLen {
// 					return s
// 				}
// 				return s[:maxLen-3] + "..."
// 			}

// 			log.Printf("Found %d servers with IP %s:", len(serverGroup), ip)
// 			for _, server := range serverGroup {
// 				log.Printf("  Name: %-50s | IP: %-21s | Bots: %2d | Players: %2d/%2d | Map: %-20s | Tags: %s",
// 					truncateString(server.Name, 50),
// 					server.Addr,
// 					server.Bots,
// 					server.Players,
// 					server.MaxPlayers,
// 					server.Map,
// 					server.GameType)
// 			}
// 			log.Printf("---")
// 		}
// 	}

// 	log.Printf("Finished analyzing duplicates. Found %d IPs with multiple servers", duplicateCount)
// }

func (c *SteamClient) FetchServers() ([]models.Server, error) {
	var allServers []models.Server
	maxRetries := 3
	// Try different regions to split up the results and avoid the 10k limit
	regions := []string{
		"\\region\\0", // US East
		"\\region\\1", // US West
		"\\region\\2", // South America
		"\\region\\3", // Europe
		"\\region\\4", // Asia
		"\\region\\5", // Australia
		"\\region\\6", // Middle East
		"\\region\\7", // Africa
	}

	for _, region := range regions {
		for i := 0; i < maxRetries; i++ {
			// Basic filter for dedicated servers and specific region
			url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&filter=\\appid\\730\\dedicated\\1%s&limit=10000", c.apiKey, region)

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Region %s, Attempt %d: HTTP error: %v", region, i+1, err)
				continue
			}

			log.Printf("Region %s, Attempt %d: Status Code: %d", region, i+1, resp.StatusCode)

			rawBody, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				log.Printf("Region %s, Attempt %d: Body read error: %v", region, i+1, err)
				continue
			}

			var result struct {
				Response struct {
					Servers []models.Server `json:"servers"`
				} `json:"response"`
			}

			if err := json.Unmarshal(rawBody, &result); err != nil {
				log.Printf("Region %s, Attempt %d: JSON parse error: %v", region, i+1, err)
				continue
			}

			// Filter and process servers
			for _, server := range result.Response.Servers {
				// Skip Valve official servers
				if strings.HasPrefix(server.Name, "Valve Counter-Strike") {
					continue
				}

				allServers = append(allServers, server)
			}

			// If we got servers, break the retry loop for this region
			if len(result.Response.Servers) > 0 {
				break
			}

			time.Sleep(time.Second * 2)
		}
	}

	log.Printf("Total servers found across all regions: %d", len(allServers))

	if len(allServers) > 10 {
		return allServers, nil
	}

	return nil, fmt.Errorf("failed to fetch sufficient servers across all regions")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
