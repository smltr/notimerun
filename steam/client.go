package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"frfr/models"
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
	url := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&filter=\\appid\\730", c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var result struct {
		Response struct {
			Servers []models.Server `json:"servers"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Response.Servers, nil
}
