package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/handler"
	"github.com/elazarl/goproxy"
)

var brainAddress string = os.Getenv("PISEC_BRAIN_ADDR")
var detailsEndpoint string = "/api/v1/indicators/details"
var indicatorsEndpoint string = "/api/v1/indicators"

var repo *cache.RedisRepository
var urlHandler *handler.PisecHandler

func main() {

	//setup the REDIS cache
	repo = cache.NewRedisClient()

	server := &handler.Server{
		BaseAddress:        brainAddress,
		IndicatorsEndpoint: indicatorsEndpoint,
		DetailsEndpoint:    detailsEndpoint,
	}

	//setup the filter
	urlHandler = handler.NewUrlHandler(repo, server)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(urlHandler.IsMalwareRequestHttp()).DoFunc(handler.GetPiSecPage)
	proxy.OnRequest(urlHandler.IsMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}
