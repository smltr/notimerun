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

func (c *SteamClient) FetchServers() ([]models.Server, error) {
	var servers []models.Server
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&filter=\\appid\\730\\dedicated\\1&limit=20000", c.apiKey)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Attempt %d: HTTP error: %v", i+1, err)
			continue
		}

		// Read and log response status
		log.Printf("Attempt %d: Status Code: %d", i+1, resp.StatusCode)

		rawBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			log.Printf("Attempt %d: Body read error: %v", i+1, err)
			continue
		}

		// Parse response
		var result struct {
			Response struct {
				Servers []models.Server `json:"servers"`
			} `json:"response"`
		}

		if err := json.Unmarshal(rawBody, &result); err != nil {
			log.Printf("Attempt %d: JSON parse error: %v", i+1, err)
			// Log the first 1000 characters of response
			log.Printf("Response preview: %s", string(rawBody[:min(1000, len(rawBody))]))
			continue
		}

		// Filter servers
		filteredServers := []models.Server{}
		for _, server := range result.Response.Servers {
			if !strings.HasPrefix(server.Name, "Valve Counter-Strike") {
				filteredServers = append(filteredServers, server)
			}
		}

		log.Printf("Attempt %d: Found %d total servers, %d after filtering",
			i+1, len(result.Response.Servers), len(filteredServers))

		if len(filteredServers) > 10 {
			return filteredServers, nil
		}

		// If we got less than 10 servers, wait and retry
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
