package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	appCache "github.com/NTCults/services-test/main-service/cache"
	"github.com/NTCults/services-test/main-service/config"
	"github.com/NTCults/services-test/main-service/models"

	"sync"
	"time"

	"github.com/NTCults/services-test/utils"
	"github.com/julienschmidt/httprouter"
)

const (
	campaigns = "campaigns"
	stats     = "stats"
	tags      = "tags"
)

var services = map[string]string{
	campaigns: "http://campaigns:8090/campaigns/",
	stats:     "http://stat:8070/stat/",
	tags:      "http://tags:8060/tags/",
}

func init() {
	servicesArray := []string{}
	for k := range services {
		servicesArray = append(servicesArray, k)
	}
	appCache.Cache.InitCache(servicesArray)
}

func main() {
	router := httprouter.New()
	router.GET("/info/:account", infoHandler)
	config := config.Config

	server := &http.Server{
		Addr:         config.Port,
		ReadTimeout:  config.ReadTimeout * time.Second,
		WriteTimeout: config.WriteTimeout * time.Second,
		Handler:      router,
	}
	fmt.Printf("Running on port %s\n", config.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accName := p.ByName("account")
	var wg sync.WaitGroup
	jsonResponses := make(chan models.ServiceResponse)
	for key, url := range services {
		wg.Add(1)
		go makeRequest(url, accName, key, &wg, jsonResponses)
	}
	collectedData := new(models.CollectedData)

	go func() {
		wg.Add(3)
		for response := range jsonResponses {
			switch response.ServiceName {
			case campaigns:
				var campArray []models.Campaign
				if err := json.Unmarshal(response.Data, &campArray); err != nil {
					fmt.Println(err)
				}
				collectedData.Campaigns = campArray
				wg.Done()
			case stats:
				var statsArray []models.Stat
				if err := json.Unmarshal(response.Data, &statsArray); err != nil {
					fmt.Println(err)
				}
				collectedData.Stats = statsArray
				wg.Done()
			case tags:
				var tagsArray []models.Tag
				if err := json.Unmarshal(response.Data, &tagsArray); err != nil {
					fmt.Println(err)
				}
				collectedData.Tags = tagsArray
				wg.Done()
			default:
				wg.Done()
			}
		}
	}()
	wg.Wait()
	data := collectedData.Aggregate()
	utils.ResponseJSON(w, http.StatusOK, data)
}

func makeRequest(url string, ID string, serviceName string, wg *sync.WaitGroup, ch chan<- models.ServiceResponse) {
	defer wg.Done()
	res, err := http.Get(url + ID)
	if err != nil {
		fmt.Println(err)
		ch <- models.ServiceResponse{serviceName, []byte{}, err}
		return
	}

	if res.StatusCode == http.StatusTooManyRequests {
		ch <- models.ServiceResponse{serviceName, appCache.Cache.Get(serviceName, ID), nil}
		fmt.Printf("Rate limit for %s has been exceeded.\n", url)
		fmt.Println("Data has been fetched from cache.")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		ch <- models.ServiceResponse{serviceName, []byte{}, err}
	}
	appCache.Cache.Set(serviceName, ID, body)
	ch <- models.ServiceResponse{serviceName, body, nil}
}
