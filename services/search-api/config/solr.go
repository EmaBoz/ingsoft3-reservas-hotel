package config

import (
	"fmt"
	"os"
)

// SolrConfig contiene la configuración de Solr
type SolrConfig struct {
	BaseURL string
	Core    string
}

// NewSolrConfig crea una nueva configuración de Solr desde variables de entorno
func NewSolrConfig() *SolrConfig {
	baseURL := os.Getenv("SOLR_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8983/solr"
	}

	core := os.Getenv("SOLR_CORE")
	if core == "" {
		core = "rooms-core"
	}

	return &SolrConfig{
		BaseURL: baseURL,
		Core:    core,
	}
}

// GetCoreURL retorna la URL completa del core de Solr
func (c *SolrConfig) GetCoreURL() string {
	return fmt.Sprintf("%s/%s", c.BaseURL, c.Core)
}

// GetSelectURL retorna la URL para consultas (select)
func (c *SolrConfig) GetSelectURL() string {
	return fmt.Sprintf("%s/%s/select", c.BaseURL, c.Core)
}

// GetUpdateURL retorna la URL para actualizaciones (indexar/eliminar)
func (c *SolrConfig) GetUpdateURL() string {
	return fmt.Sprintf("%s/%s/update/json", c.BaseURL, c.Core)
}
