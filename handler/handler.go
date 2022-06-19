package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Ringloop/pisec/cache"
	"github.com/Ringloop/pisec/filter"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/elazarl/goproxy"
)

type PisecHandler struct {
	urlFilter *filter.PisecUrlFilter
}

type Server struct {
	BaseAddress        string
	IndicatorsEndpoint string
	DetailsEndpoint    string
}

func downloadBloomFilter(indicatorsEndpoint string) *bloom.BloomFilter {

	var filter *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)

	//download the bloom filter from server
	res, err := http.Get(indicatorsEndpoint)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	jsonRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = filter.UnmarshalJSON(jsonRes)
	if err != nil {
		panic(err)
	}

	return filter
}

func NewUrlHandler(repo *cache.RedisRepository, server *Server) *PisecHandler {

	bloomFilter := downloadBloomFilter(server.BaseAddress + server.IndicatorsEndpoint)
	urlFilter := filter.NewPisecUrlFilter(repo, bloomFilter, server.BaseAddress+server.IndicatorsEndpoint)

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
