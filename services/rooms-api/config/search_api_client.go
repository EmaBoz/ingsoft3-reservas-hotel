package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"rooms-api/domain"
)

// SearchAPIClient es el cliente HTTP para search-api
type SearchAPIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewSearchAPIClient crea un nuevo cliente para search-api
func NewSearchAPIClient() *SearchAPIClient {
	baseURL := os.Getenv("SEARCH_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://search-api:8083"
	}

	return &SearchAPIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchRoomsResponse representa la respuesta de search-api
type SearchRoomsResponse struct {
	Page    int                `json:"page"`
	Limit   int                `json:"limit"`
	Total   int64              `json:"total"`
	Results []SearchRoomResult `json:"results"`
}

// SearchRoomResult representa un resultado individual de búsqueda
type SearchRoomResult struct {
	ID          []string  `json:"id"`
	Number      []int     `json:"number"`
	Type        []string  `json:"type"`
	Status      []string  `json:"status"`
	Price       []float64 `json:"price"`
	Description []string  `json:"description,omitempty"`
	Capacity    []int     `json:"capacity"`
	Floor       []int     `json:"floor"`
	HasWifi     []bool    `json:"has_wifi"`
	HasAC       []bool    `json:"has_ac"`
	HasTV       []bool    `json:"has_tv"`
	HasMinibar  []bool    `json:"has_minibar"`
	CreatedAt   []string  `json:"created_at"`
	UpdatedAt   []string  `json:"updated_at"`
}

// SearchRooms realiza una búsqueda en search-api
func (c *SearchAPIClient) SearchRooms(filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	// Construir query parameters
	params := url.Values{}

	if filter.Type != nil {
		params.Add("type", string(*filter.Type))
	}
	if filter.Status != nil {
		params.Add("status", string(*filter.Status))
	}
	if filter.Floor != nil {
		params.Add("floor", strconv.Itoa(*filter.Floor))
	}
	if filter.MinPrice != nil {
		params.Add("min_price", fmt.Sprintf("%.2f", *filter.MinPrice))
	}
	if filter.MaxPrice != nil {
		params.Add("max_price", fmt.Sprintf("%.2f", *filter.MaxPrice))
	}
	if filter.HasWifi != nil {
		params.Add("has_wifi", strconv.FormatBool(*filter.HasWifi))
	}
	if filter.HasAC != nil {
		params.Add("has_ac", strconv.FormatBool(*filter.HasAC))
	}
	if filter.HasTV != nil {
		params.Add("has_tv", strconv.FormatBool(*filter.HasTV))
	}
	if filter.HasMinibar != nil {
		params.Add("has_minibar", strconv.FormatBool(*filter.HasMinibar))
	}

	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Construir URL
	searchURL := fmt.Sprintf("%s/api/search/rooms?%s", c.BaseURL, params.Encode())

	// Realizar request
	resp, err := c.HTTPClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request search-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search-api returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decodificar respuesta
	var searchResp SearchRoomsResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search-api response: %w", err)
	}

	// Convertir resultados de Solr a RoomResponse
	rooms := make([]domain.RoomResponse, len(searchResp.Results))
	for i, result := range searchResp.Results {
		// Extraer primer elemento de cada array (Solr devuelve arrays)
		idStr := ""
		if len(result.ID) > 0 {
			idStr = result.ID[0]
		}

		id, _ := strconv.ParseUint(idStr, 10, 32)
		number := ""
		if len(result.Number) > 0 {
			number = strconv.Itoa(result.Number[0])
		}

		rooms[i] = domain.RoomResponse{
			ID:          uint(id),
			Number:      number,
			Type:        domain.RoomType(getFirstString(result.Type)),
			Status:      domain.RoomStatus(getFirstString(result.Status)),
			Price:       getFirstFloat64(result.Price),
			Description: getFirstString(result.Description),
			Capacity:    getFirstInt(result.Capacity),
			Floor:       getFirstInt(result.Floor),
			HasWifi:     getFirstBool(result.HasWifi),
			HasAC:       getFirstBool(result.HasAC),
			HasTV:       getFirstBool(result.HasTV),
			HasMinibar:  getFirstBool(result.HasMinibar),
		}
	}

	return &domain.RoomListResponse{
		Rooms: rooms,
		Total: searchResp.Total,
		Page:  searchResp.Page,
		Limit: searchResp.Limit,
	}, nil
}

// HealthCheck verifica si search-api está disponible
func (c *SearchAPIClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return fmt.Errorf("search-api health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search-api health check returned status %d", resp.StatusCode)
	}

	return nil
}

// Funciones helper para extraer primer elemento de arrays de Solr
func getFirstString(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func getFirstInt(arr []int) int {
	if len(arr) > 0 {
		return arr[0]
	}
	return 0
}

func getFirstFloat64(arr []float64) float64 {
	if len(arr) > 0 {
		return arr[0]
	}
	return 0.0
}

func getFirstBool(arr []bool) bool {
	if len(arr) > 0 {
		return arr[0]
	}
	return false
}
