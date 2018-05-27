package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Campaign struct represents main service entity
type Campaign struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Get returns slice of Campaign objects
func (c *Campaign) Get(accountID string, db *sql.DB) ([]Campaign, error) {
	campaigns := make([]Campaign, 0)

	query := fmt.Sprintf("SELECT id, title FROM campaigns WHERE account='%s'", accountID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, errors.New("something wrong")
	}

	for rows.Next() {
		campaign := new(Campaign)
		err = rows.Scan(&campaign.ID, &campaign.Title)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Something wrong")
		}
		campaigns = append(campaigns, *campaign)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Something wrong")
	}

	return campaigns, err
}
