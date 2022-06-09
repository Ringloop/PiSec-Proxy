package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	proxy.OnRequest(isMalwareRequestHttp()).DoFunc(GetPiSecPage)
	proxy.OnRequest(isMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}

func GetPiSecPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Blocked By PiSec with <3 !")
}

var IsConnectToMalware goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	fmt.Println("connecting...")
	return goproxy.MitmConnect, host
}

func isMalwareRequestHttp() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTP")
		res, err := urlFilter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

func isMalwareRequestHttps() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTPS")
		res, err := urlFilter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}
