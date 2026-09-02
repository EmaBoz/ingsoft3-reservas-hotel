package config

import (
	"os"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/karlseguin/ccache/v3"
)

// CacheConfig contiene la configuración de caché
type CacheConfig struct {
	LocalTTL        time.Duration
	DistributedTTL  time.Duration
	MemcachedHost   string
	MemcachedPort   string
	LocalMaxSize    int64
	LocalItemsToPrune uint32
}

// NewCacheConfig crea una nueva configuración de caché
func NewCacheConfig() *CacheConfig {
	localTTL := 60 * time.Second
	if envTTL := os.Getenv("LOCAL_CACHE_TTL_SECONDS"); envTTL != "" {
		if seconds, err := strconv.Atoi(envTTL); err == nil {
			localTTL = time.Duration(seconds) * time.Second
		}
	}

	distributedTTL := 300 * time.Second
	if envTTL := os.Getenv("DISTRIBUTED_CACHE_TTL_SECONDS"); envTTL != "" {
		if seconds, err := strconv.Atoi(envTTL); err == nil {
			distributedTTL = time.Duration(seconds) * time.Second
		}
	}

	memcachedHost := os.Getenv("MEMCACHED_HOST")
	if memcachedHost == "" {
		memcachedHost = "localhost"
	}

	memcachedPort := os.Getenv("MEMCACHED_PORT")
	if memcachedPort == "" {
		memcachedPort = "11211"
	}

	return &CacheConfig{
		LocalTTL:          localTTL,
		DistributedTTL:    distributedTTL,
		MemcachedHost:     memcachedHost,
		MemcachedPort:     memcachedPort,
		LocalMaxSize:      500 * 1024 * 1024, // 500MB
		LocalItemsToPrune: 50,
	}
}

// NewLocalCache crea una nueva instancia de CCache
func (c *CacheConfig) NewLocalCache() *ccache.Cache[[]byte] {
	config := ccache.Configure[[]byte]()
	config.MaxSize(c.LocalMaxSize)
	config.ItemsToPrune(c.LocalItemsToPrune)

	return ccache.New(config)
}

// NewMemcachedClient crea un nuevo cliente de Memcached
func (c *CacheConfig) NewMemcachedClient() *memcache.Client {
	server := c.MemcachedHost + ":" + c.MemcachedPort
	return memcache.New(server)
}

// GetMemcachedAddress retorna la dirección completa de Memcached
func (c *CacheConfig) GetMemcachedAddress() string {
	return c.MemcachedHost + ":" + c.MemcachedPort
}
