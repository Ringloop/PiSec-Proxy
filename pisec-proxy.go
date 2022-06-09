package main

import (
	"log"
	"net/http"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/filter"
	"github.com/Ringloop/pisec/handler"
	"github.com/elazarl/goproxy"
)

var urlFilter *filter.PisecUrlFilter
var repo *cache.RedisRepository

func main() {

	//setup the REDIS cache
	repo = cache.NewRedisClient()

	//setup the filter
	urlFilter = filter.NewPisecUrlFilter(repo)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(handler.IsMalwareRequestHttp()).DoFunc(handler.GetPiSecPage)
	proxy.OnRequest(handler.IsMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}
