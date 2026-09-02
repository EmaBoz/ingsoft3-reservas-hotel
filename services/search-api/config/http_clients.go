package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"search-api/domain"
)

// RoomsAPIClient es el cliente HTTP para rooms-api
type RoomsAPIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewRoomsAPIClient crea un nuevo cliente para rooms-api
func NewRoomsAPIClient() *RoomsAPIClient {
	baseURL := os.Getenv("ROOMS_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	return &RoomsAPIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRoomByID obtiene una habitación por ID desde rooms-api
func (c *RoomsAPIClient) GetRoomByID(id uint) (*domain.Room, error) {
	url := fmt.Sprintf("%s/api/v1/rooms/%d", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request rooms-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("room %d not found", id)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rooms-api returned status %d: %s", resp.StatusCode, string(body))
	}

	// rooms-api devuelve directamente el objeto room sin wrapper
	var room domain.Room
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return nil, fmt.Errorf("failed to decode rooms-api response: %w", err)
	}

	return &room, nil
}

// HealthCheck verifica si rooms-api está disponible
func (c *RoomsAPIClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return fmt.Errorf("rooms-api health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rooms-api health check returned status %d", resp.StatusCode)
	}

	return nil
}
