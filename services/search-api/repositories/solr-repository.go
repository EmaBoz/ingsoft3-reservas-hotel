package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"search-api/config"
	"search-api/domain"
)

// SolrRepository maneja las operaciones con Solr
type SolrRepository struct {
	config     *config.SolrConfig
	httpClient *http.Client
}

// NewSolrRepository crea un nuevo repositorio de Solr
func NewSolrRepository(cfg *config.SolrConfig) *SolrRepository {
	return &SolrRepository{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SolrResponse representa la respuesta de Solr para búsquedas
type SolrResponse struct {
	Response struct {
		NumFound int64                     `json:"numFound"`
		Start    int                       `json:"start"`
		Docs     []domain.SolrRoomDocument `json:"docs"`
	} `json:"response"`
}

// SolrUpdateResponse representa la respuesta de Solr para updates
type SolrUpdateResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
}

// Search realiza una búsqueda en Solr
func (r *SolrRepository) Search(req *domain.SearchRoomsRequest) (*domain.SearchRoomsResponse, error) {
	// Construir query Solr
	params := r.buildSolrQuery(req)

	// Realizar request
	url := r.config.GetSelectURL() + "?" + params
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Solr returned status %d: %s", resp.StatusCode, string(body))
	}

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		return nil, fmt.Errorf("failed to decode Solr response: %w", err)
	}

	// Mapear a respuesta
	return &domain.SearchRoomsResponse{
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   solrResp.Response.NumFound,
		Results: solrResp.Response.Docs,
	}, nil
}

// buildSolrQuery construye los parámetros de query para Solr
func (r *SolrRepository) buildSolrQuery(req *domain.SearchRoomsRequest) string {
	params := url.Values{}

	// Query principal
	if req.Q != "" {
		// Búsqueda en múltiples campos
		params.Add("q", fmt.Sprintf("number:*%s* OR type:*%s* OR description:*%s*", req.Q, req.Q, req.Q))
	} else {
		params.Add("q", "*:*")
	}

	// Filtros (fq)
	filters := []string{}

	if req.Type != "" {
		filters = append(filters, fmt.Sprintf("type:%s", req.Type))
	}

	if req.Status != "" {
		filters = append(filters, fmt.Sprintf("status:%s", req.Status))
	}

	if req.Floor != nil {
		filters = append(filters, fmt.Sprintf("floor:%d", *req.Floor))
	}

	if req.MinPrice != nil && req.MaxPrice != nil {
		filters = append(filters, fmt.Sprintf("price:[%.0f TO %.0f]", *req.MinPrice, *req.MaxPrice))
	} else if req.MinPrice != nil {
		filters = append(filters, fmt.Sprintf("price:[%.0f TO *]", *req.MinPrice))
	} else if req.MaxPrice != nil {
		filters = append(filters, fmt.Sprintf("price:[* TO %.0f]", *req.MaxPrice))
	}
	if req.HasWifi != nil {
		filters = append(filters, fmt.Sprintf("has_wifi:%t", *req.HasWifi))
	}

	if req.HasAC != nil {
		filters = append(filters, fmt.Sprintf("has_ac:%t", *req.HasAC))
	}

	if req.HasTV != nil {
		filters = append(filters, fmt.Sprintf("has_tv:%t", *req.HasTV))
	}

	if req.HasMinibar != nil {
		filters = append(filters, fmt.Sprintf("has_minibar:%t", *req.HasMinibar))
	}

	if len(filters) > 0 {
		params.Add("fq", strings.Join(filters, " AND "))
	}

	// Ordenamiento
	if req.Sort != "" {
		sort := r.parseSortParam(req.Sort)
		if sort != "" {
			params.Add("sort", sort)
		}
	} else {
		params.Add("sort", "id desc") // Default
	}

	// Paginación
	start := (req.Page - 1) * req.Limit
	params.Add("start", strconv.Itoa(start))
	params.Add("rows", strconv.Itoa(req.Limit))

	// Formato
	params.Add("wt", "json")

	// Convertir a query string con encoding correcto
	return params.Encode()
}

// parseSortParam parsea el parámetro de sort (ej: "price", "-price")
func (r *SolrRepository) parseSortParam(sort string) string {
	if sort == "" {
		return ""
	}

	order := "asc"
	field := sort

	if strings.HasPrefix(sort, "-") {
		order = "desc"
		field = sort[1:]
	}

	// Validar campos permitidos
	validFields := map[string]bool{
		"price": true, "floor": true, "capacity": true,
		"number": true, "id": true,
	}

	if !validFields[field] {
		return ""
	}

	return fmt.Sprintf("%s %s", field, order)
}

// IndexDocument indexa un documento en Solr
func (r *SolrRepository) IndexDocument(doc *domain.SolrRoomWrite) error {
	// Enviar el documento directamente como array de un elemento
	payload := []interface{}{doc}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	url := r.config.GetUpdateURL() + "?commit=true"
	resp, err := r.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to index document in Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteDocument elimina un documento de Solr por ID
func (r *SolrRepository) DeleteDocument(id string) error {
	payload := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal delete request: %w", err)
	}

	url := r.config.GetUpdateURL() + "?commit=true"
	resp, err := r.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to delete document from Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// HealthCheck verifica la salud de Solr
func (r *SolrRepository) HealthCheck() error {
	url := r.config.GetCoreURL() + "/admin/ping"
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("Solr health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Solr health check returned status %d", resp.StatusCode)
	}

	return nil
}
