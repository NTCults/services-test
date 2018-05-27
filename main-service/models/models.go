package models

import (
	"encoding/json"
	"fmt"

	"github.com/NTCults/services-test/main-service/config"
)

// ServiceResponse represent data fetched from outer service
type ServiceResponse struct {
	ServiceName string
	Data        []byte
	Err         error
}

// CollectedData represent data collected from outer services
type CollectedData struct {
	Campaigns []campaign
	Stats     []stat
	Tags      []tag
}

// HandleResponse collects data from ServiceResponse
func (c *CollectedData) HandleResponse(sr ServiceResponse) error {
	if sr.Err != nil {
		return sr.Err
	}
	switch sr.ServiceName {
	case config.CampaignsServiceName:
		var campArray []campaign
		if err := json.Unmarshal(sr.Data, &campArray); err != nil {
			fmt.Println(err)
			return err
		}
		c.Campaigns = campArray
	case config.StatsServiceName:
		var statsArray []stat
		if err := json.Unmarshal(sr.Data, &statsArray); err != nil {
			fmt.Println(err)
			return err
		}
		c.Stats = statsArray
	case config.TagsServiceName:
		var tagsArray []tag
		if err := json.Unmarshal(sr.Data, &tagsArray); err != nil {
			fmt.Println(err)
			return err
		}
		c.Tags = tagsArray
	}
	return nil
}

// Aggregate aggregates data collected from outer services
func (c *CollectedData) Aggregate() *[]Info {
	var infoArray []Info
	for _, camp := range c.Campaigns {
		var info Info
		campID := camp.ID

		info.ID = campID
		info.Title = camp.Title

		for _, stat := range c.Stats {
			if *stat.CampaignID == info.ID {
				stat.CampaignID = nil
				info.Stat = append(info.Stat, stat)
			}
		}

		for _, tag := range c.Tags {
			if *tag.CampaignID == info.ID {
				tag.CampaignID = nil
				info.Tags = append(info.Tags, tag)
			}
		}
		infoArray = append(infoArray, info)
	}
	return &infoArray
}

// Campaign represent data collected from campaigns service
type campaign struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Stat represent data collected from stat service
type stat struct {
	CampaignID *int   `json:"campaign_id,omitempty"`
	Date       string `json:"date"`
	Shows      int    `json:"shows"`
	Clicks     int    `json:"clicks"`
	Costs      int    `json:"costs"`
}

// Tag represent data collected from tag service
type tag struct {
	CampaignID *int     `json:"campaign_id,omitempty"`
	Tags       []string `json:"tags"`
}

// Info represents unit of main services output
type Info struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Stat  []stat `json:"stat"`
	Tags  []tag  `json:"tags"`
}
