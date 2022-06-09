package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Ringloop/pisec/filter"
	"github.com/elazarl/goproxy"
)

func GetPiSecPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Blocked By PiSec with <3 !")
}

var IsConnectToMalware goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	fmt.Println("connecting...")
	return goproxy.MitmConnect, host
}

func IsMalwareRequestHttp() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTP")
		res, err := filter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

func IsMalwareRequestHttps() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTPS")
		res, err := filter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}
