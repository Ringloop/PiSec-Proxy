package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(IsMalware()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if h, _, _ := time.Now().Clock(); h >= 8 && h <= 24 {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Blocked By PiSec with <3 !")
			}
			return r, nil
		})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func IsMalware() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.URL.Host == "news.ycombinator.com"
	}
}
