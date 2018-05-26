package main

import (
	"net/http"
	"services-test/tags/utils"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

type tags struct {
	CampaignID int      `json:"campaign_id"`
	Tags       []string `json:""tags`
}

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
	router.GET("/tags/:account", tagsHandler)
	http.ListenAndServe(":8060", limit(router))
}

func tagsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	testA := tags{12, []string{"test", "test3"}}
	testB := tags{12, []string{"test", "test5"}}

	var data []tags

	data = append(data, testA)
	data = append(data, testB)

	// accName := p.ByName("account")
	utils.ResponseJSON(w, http.StatusOK, data)
}
