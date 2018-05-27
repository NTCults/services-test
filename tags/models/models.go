package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Tags struct {
	CampaignID int      `json:"campaign_id"`
	Tags       []string `json:"tags"`
}

type tag struct {
	tag        string
	campaignID int
}

func aggregateTags(tags map[int][]string) []Tags {
	result := []Tags{}
	for k, v := range tags {
		var tagsByCampaign Tags
		tagsByCampaign.CampaignID = k
		tagsByCampaign.Tags = v
		result = append(result, tagsByCampaign)
	}
	return result
}

func (c *Tags) Get(accountID string, db *sql.DB) ([]Tags, error) {
	tags := make(map[int][]string)
	query := fmt.Sprintf(`SELECT campaign_id, tag FROM tags WHERE account='%s'`, accountID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Something wrong")
	}

	for rows.Next() {
		tag := new(tag)
		err = rows.Scan(&tag.campaignID, &tag.tag)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Something wrong")
		}
		tags[tag.campaignID] = append(tags[tag.campaignID], tag.tag)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Something wrong")
	}
	result := aggregateTags(tags)
	return result, err
}
