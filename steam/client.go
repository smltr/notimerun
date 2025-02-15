package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"notimerun/models"
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
	url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&filter=\\appid\\730\\dedicated\\1&limit=20000", c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the raw response body
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result struct {
		Response struct {
			Servers []models.Server `json:"servers"`
		} `json:"response"`
	}

	if err := json.Unmarshal(rawBody, &result); err != nil {
		return nil, err
	}

	// Filter out Valve official servers
	var filteredServers []models.Server
	for _, server := range result.Response.Servers {
		// Skip servers with IP starting with "100.66"
		if !strings.HasPrefix(server.Name, "Valve Counter-Strike") {
			filteredServers = append(filteredServers, server)
		}
	}

	// Write first 100 filtered servers to file for debugging
	if len(filteredServers) > 100 {
		debugServers := filteredServers[:100]
		debugJson, err := json.MarshalIndent(debugServers, "", "    ")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile("tmp/servers.json", debugJson, 0644)
		if err != nil {
			return nil, err
		}
	}

	return filteredServers, nil
}
