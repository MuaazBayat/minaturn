package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// APIClient handles communication with Django API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient creates a new client for Django API
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAllQueues fetches all queues and their entries from Django API
func (c *APIClient) GetAllQueues() (*APIResponse, error) {
	url := fmt.Sprintf("%s/queues/all/", c.BaseURL)
	
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch queues: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp, nil
}

// parseTime safely parses ISO timestamp strings, handling null values
func parseTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}