package main

import (
	"log"
	"net/http"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/handler"
	"github.com/elazarl/goproxy"
)

var repo *cache.RedisRepository
var urlHandler *handler.PisecHandler

func main() {

	//setup the REDIS cache
	repo = cache.NewRedisClient()

	//setup the filter
	urlHandler = handler.NewUrlHandler(repo)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(urlHandler.IsMalwareRequestHttp()).DoFunc(handler.GetPiSecPage)
	proxy.OnRequest(urlHandler.IsMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}
