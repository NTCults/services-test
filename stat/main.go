package main

import (
	"fmt"
	"net/http"
	"github.com/NTCults/services-test/stat/db"
	"github.com/NTCults/services-test/stat/models"

	"github.com/NTCults/services-test/utils"

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
	router.GET("/stat/:account", statHandler)
	http.ListenAndServe(":8070", limit(router))
}

func statHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	account := p.ByName("account")
	stat := new(models.Stat)
	result, err := stat.Get(account, db.Connect)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusOK, result)
}
