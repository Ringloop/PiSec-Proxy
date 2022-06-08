package main

import (
	"log"
	"net/http"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/url_filter"
	"github.com/elazarl/goproxy"
)

var urlFilter *url_filter.PisecUrlFilter
var repo *cache.RedisRepository

func main() {

	//setup the REDIS cache
	repo = cache.NewRedisClient()

	urlFilter = url_filter.NewPisecUrlFilter(repo)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(urlFilter.IsMalwareRequestHttp()).DoFunc(GetPiSecPage)
	proxy.OnRequest(urlFilter.IsMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}

func GetPiSecPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Blocked By PiSec with <3 !")
}
