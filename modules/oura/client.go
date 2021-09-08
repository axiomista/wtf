package oura

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// var (
// 	accessToken string // Required config value
// 	dayOffset   int
// )

// Client represents the data required to connect to the Twitter API
type Client struct {
	apiBase     string
	accessToken string
	start       string
	end         string
	httpClient  *http.Client
}

// NewClient creates and returns a new Twitter client
func NewClient(settings *Settings) *Client {
	offSetDays := settings.days * -1
	start := time.Now().AddDate(0, 0, offSetDays).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 0).Format("2006-01-02") // Today

	return &Client{
		apiBase:     "https://api.ouraring.com/v1/",
		accessToken: settings.accessToken,
		start:       start,
		end:         end,
		httpClient:  http.DefaultClient,
	}
}

// getUserInfo returns personal data for a user
func (client *Client) getUserInfo() (*UserInfo, error) {
	var userInfo UserInfo
	apiURL := fmt.Sprintf("%suserinfo", client.apiBase)
	data, err := get(client.httpClient, apiURL, client.accessToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

// getSleeps returns a slice of sleep data for a user
func (client *Client) getSleeps() ([]Sleep, error) {
	var sleepDays SleepSummaries
	apiURL := fmt.Sprintf("%ssleep?start=%s&end=%s", client.apiBase, client.start, client.end)
	data, err := get(client.httpClient, apiURL, client.accessToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &sleepDays)
	if err != nil {
		return nil, err
	}
	return sleepDays.SleepDays, nil
}

// getReadiness returns a slice of readiness data for a user
func (client *Client) getReadiness() ([]Readiness, error) {
	var readinessDays ReadinessSummaries
	apiURL := fmt.Sprintf("%sreadiness?start=%s&end=%s", client.apiBase, client.start, client.end)
	data, err := get(client.httpClient, apiURL, client.accessToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &readinessDays)
	if err != nil {
		return nil, err
	}
	return readinessDays.ReadinessDays, nil
}

// getActivity returns a slice of activity data for a user
func (client *Client) getActivity() ([]Activity, error) {
	var activityDays ActivitySummaries
	apiURL := fmt.Sprintf("%sactivity?start=%s&end=%s", client.apiBase, client.start, client.end)
	data, err := get(client.httpClient, apiURL, client.accessToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &activityDays)
	if err != nil {
		return nil, err
	}
	return activityDays.ActivityDays, nil
}

// getBedtimes returns a slice of ideal bedtime data for a user
func (client *Client) getBedtimes() ([]Bedtime, error) {
	var idealBedtimes Bedtimes
	apiURL := fmt.Sprintf("%sbedtime?start=%s&end=%s", client.apiBase, client.start, client.end)
	data, err := get(client.httpClient, apiURL, client.accessToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &idealBedtimes)
	if err != nil {
		return nil, err
	}
	return idealBedtimes.IdealBedtimes, nil
}

func get(httpClient *http.Client, apiURL string, accessToken string) ([]byte, error) {
	// Build request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send request and ensure success
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Request to %s returned status code %d", apiURL, response.StatusCode)
	}

	// Reads response body to []byte
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
