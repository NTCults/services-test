package main

import (
	"net/http"
	"services-test/campaigns/utils"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

type campaign struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type campaigns []campaign

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
	router.GET("/campaigns/:account", campaignsHandler)
	http.ListenAndServe(":8090", limit(router))
}

func campaignsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	testA := campaign{12, "test"}
	testB := campaign{12, "test"}
	var data campaigns
	data = append(data, testA)
	data = append(data, testB)
	utils.ResponseJSON(w, http.StatusOK, data)
}
