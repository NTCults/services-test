package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

var Config ServerConfig

// Constants useing to indentify outer services
const (
	CampaignsServiceName = "campaigns"
	StatsServiceName     = "stats"
	TagsServiceName      = "tags"
)

type ServerConfig struct {
	Port                string        `json:"port"`
	ReadTimeout         time.Duration `json:"write_timeout"`
	WriteTimeout        time.Duration `json:"read_timeout"`
	CacheExpirationTime int           `json:"cache_expiration_time"`
	CachePurgesTime     int           `json:"cache_purges_time"`
}

func init() {
	config := new(ServerConfig)
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		log.Fatal("Configuration error:", err)
	}
	Config = *config
}
