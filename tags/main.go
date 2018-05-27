package main

import (
	"fmt"
	"net/http"
	"services-test/tags/db"
	"services-test/tags/models"

	"github.com/NTCults/services-test/campaigns/utils"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

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
	db.ConnectToDB()
	router.GET("/tags/:account", tagsHandler)
	http.ListenAndServe(":8060", limit(router))
}

func tagsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	account := p.ByName("account")
	tags := new(models.Tags)
	result, err := tags.Get(account, db.Connect)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusOK, result)
}
