package main

import (
	"net/http"
	"services-test/stat/utils"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

type stat struct {
	CampaignID int       `json:"campaign_id"`
	Date       time.Time `json:"data"`
	Shows      int       `json:"shows"`
	Clicks     int       `json:"clicks"`
	Costs      int       `json:"costs"`
}

type stats []stat

var limiter = rate.NewLimiter(1, 1)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			utils.ResponseError(w, http.StatusTooManyRequests, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := httprouter.New()
	router.GET("/stat/:account", statHandler)
	http.ListenAndServe(":8070", limit(router))
}

func statHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	testA := stat{12, time.Now(), 500, 300, 100}
	testB := stat{12, time.Now(), 500, 200, 15}

	var data stats

	data = append(data, testA)
	data = append(data, testB)

	utils.ResponseJSON(w, http.StatusOK, data)
}
