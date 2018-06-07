package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	appCache "github.com/NTCults/services-test/main-service/cache"
	"github.com/NTCults/services-test/main-service/config"
	"github.com/NTCults/services-test/main-service/models"

	"github.com/NTCults/services-test/utils"
	"github.com/julienschmidt/httprouter"
)

var services = map[string]string{
	config.CampaignsServiceName: "http://campaigns:8090/campaigns/",
	config.StatsServiceName:     "http://stat:8070/stat/",
	config.TagsServiceName:      "http://tags:8060/tags/",
}

var servicesLocal = map[string]string{
	config.CampaignsServiceName: "http://localhost:8090/campaigns/",
	config.StatsServiceName:     "http://localhost:8070/stat/",
	config.TagsServiceName:      "http://localhost:8060/tags/",
}

func init() {
	local := flag.Bool("l", false, "Use localhost url for outer services.")
	flag.Parse()
	if *local {
		services = servicesLocal
	}
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
	wg.Add(len(services))
	go func() {
		defer close(jsonResponses)
		for response := range jsonResponses {
			if err := collectedData.HandleResponse(response); err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}
	}()
	wg.Wait()

	data, err := collectedData.Aggregate()
	if err != nil {
		utils.ResponseError(w, http.StatusNotFound, err.Error())
		return
	}
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
		return
	}
	appCache.Cache.Set(serviceName, ID, body)
	ch <- models.ServiceResponse{serviceName, body, nil}
}
