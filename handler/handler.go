package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/filter"
	"github.com/elazarl/goproxy"
)

type PisecHandler struct {
	urlFilter *filter.PisecUrlFilter
}

func NewUrlHandler(repo *cache.RedisRepository) *PisecHandler {

	urlFilter := filter.NewPisecUrlFilter(repo)

	return &PisecHandler{urlFilter: urlFilter}
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

func (handler *PisecHandler) IsMalwareRequestHttp() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTP")
		res, err := handler.urlFilter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

func (handler *PisecHandler) IsMalwareRequestHttps() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTPS")
		res, err := handler.urlFilter.ShallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}
