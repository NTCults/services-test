package main

import (
	"fmt"
	"net/http"
	"services-test/campaigns/db"
	"services-test/campaigns/models"

	"github.com/NTCults/services-test/utils"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

type campaigns []models.Campaign

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
	db.ConnectToDB()
	router := httprouter.New()
	router.GET("/campaigns/:account", campaignsHandler)
	http.ListenAndServe(":8090", limit(router))
}

func campaignsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	account := p.ByName("account")
	cmp := new(models.Campaign)
	result, err := cmp.Get(account, db.Connect)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusOK, result)
}
