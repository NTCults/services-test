package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Stat struct {
	CampaignID int    `json:"campaign_id"`
	Date       string `json:"date"`
	Shows      int    `json:"shows"`
	Clicks     int    `json:"clicks"`
	Costs      int    `json:"costs"`
}

func (c *Stat) Get(accountID string, db *sql.DB) ([]Stat, error) {
	stats := make([]Stat, 0)

	query := fmt.Sprintf(`SELECT campaign_id, date, shows, clicks, costs FROM stat WHERE account='%s'`, accountID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Something wrong")
	}

	for rows.Next() {
		stat := new(Stat)
		err = rows.Scan(&stat.CampaignID, &stat.Date, &stat.Shows, &stat.Clicks, &stat.Costs)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Something wrong")
		}
		stats = append(stats, *stat)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Something wrong")
	}

	return stats, err
}
