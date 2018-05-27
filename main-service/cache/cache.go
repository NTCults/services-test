package utils

import (
	"github.com/NTCults/services-test/main-service/config"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache is main app cache instance
var Cache = new(AppCache)

// AppCache wraps map with instances of pointers to cache.Cache
type AppCache struct {
	stores map[string]*cache.Cache
}

// InitCache initializes applications cache; should be called before the server start
func (a *AppCache) InitCache(services []string) {
	expirationTime := config.Config.CacheExpirationTime
	purgesTime := config.Config.CachePurgesTime

	a.stores = make(map[string]*cache.Cache)
	for _, s := range services {
		a.stores[s] = cache.New(time.Duration(expirationTime)*time.Minute,
			time.Duration(purgesTime)*time.Minute)
	}
}

// Get is Cache.cache Get method wrapper
func (a *AppCache) Get(store string, ID string) []byte {
	var result []byte
	data, ok := a.stores[store].Get(ID)
	if ok {
		result = data.([]byte)
	}
	return result
}

// Set is Cache.cache Set method wrapper
func (a *AppCache) Set(store string, ID string, data interface{}) {
	a.stores[store].Set(ID, data, cache.DefaultExpiration)
}
