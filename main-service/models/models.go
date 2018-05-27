package models

// ServiceResponse represent data fetched from outer service
type ServiceResponse struct {
	ServiceName string
	Data        []byte
	Err         error
}

// CollectedData represent data collected from outer services
type CollectedData struct {
	Campaigns []Campaign
	Stats     []Stat
	Tags      []Tag
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
type Campaign struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Stat represent data collected from stat service
type Stat struct {
	CampaignID *int   `json:"campaign_id,omitempty"`
	Date       string `json:"date"`
	Shows      int    `json:"shows"`
	Clicks     int    `json:"clicks"`
	Costs      int    `json:"costs"`
}

// Tag represent data collected from tag service
type Tag struct {
	CampaignID *int     `json:"campaign_id,omitempty"`
	Tags       []string `json:"tags"`
}

// Info represents unit of main services output
type Info struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Stat  []Stat `json:"stat"`
	Tags  []Tag  `json:"tags"`
}
