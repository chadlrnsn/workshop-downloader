package steamcmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PublishedFileDetails struct {
	PublishedFileID string `json:"publishedfileid"`
	Result          int    `json:"result"`
	ConsumerAppID   int    `json:"consumer_app_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	PreviewURL      string `json:"preview_url"`
}

type SteamWebAPIResponse struct {
	Response struct {
		Result               int                    `json:"result"`
		ResultCount          int                    `json:"resultcount"`
		PublishedFileDetails []PublishedFileDetails `json:"publishedfiledetails"`
	} `json:"response"`
}

// FetchWorkshopMetadata queries Steam WebAPI for item metadata
func FetchWorkshopMetadata(workshopID string) (*PublishedFileDetails, error) {
	apiURL := "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"
	
	formData := url.Values{}
	formData.Set("itemcount", "1")
	formData.Set("publishedfileids[0]", workshopID)

	resp, err := http.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("metadata API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metadata API returned status code: %d", resp.StatusCode)
	}

	var apiResp SteamWebAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode metadata response: %w", err)
	}

	if apiResp.Response.ResultCount == 0 || len(apiResp.Response.PublishedFileDetails) == 0 {
		return nil, fmt.Errorf("no metadata returned for workshop ID")
	}

	details := apiResp.Response.PublishedFileDetails[0]
	if details.Result != 1 {
		return nil, fmt.Errorf("steam API returned error result code: %d (File Not Found or Private)", details.Result)
	}

	return &details, nil
}
