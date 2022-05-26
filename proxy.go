package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/elazarl/goproxy"
)

var filter *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)
var serverAddress string = os.Getenv("PISEC_BRAIN_ADDR")
var indicatorsEndpoint string = "/api/v1/indicators"
var detailsEndpoint string = "/api/v1/indicators/details"

type CheckUrlRequest struct {
	Url string `json:"url"`
}

type CheckUrlResponse struct {
	Result bool `json:"exists"`
}

func main() {

	//download the bloom filter from server
	endpoint := serverAddress + indicatorsEndpoint
	res, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	jsonRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	filter.UnmarshalJSON(jsonRes)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(IsMalwareRequestHttp()).DoFunc(GetPiSecPage)
	proxy.OnRequest(IsMalwareRequestHttps()).HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8880", proxy))
}

func GetPiSecPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Blocked By PiSec with <3 !")
}

func GetPiSecPage2(r *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Blocked By PiSec with <3 !")
}

func CheckUrlWithBrain(url string) (bool, error) {
	checkUrlTest := CheckUrlRequest{
		Url: url,
	}
	var checkUrlBuf bytes.Buffer
	err := json.NewEncoder(&checkUrlBuf).Encode(checkUrlTest)
	if err != nil {
		return false, err
	}

	endpoint := serverAddress + detailsEndpoint
	res, err := http.Post(endpoint, "application/json", &checkUrlBuf)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	jsonRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var checkUrlRes CheckUrlResponse
	json.Unmarshal(jsonRes, &checkUrlRes)
	return checkUrlRes.Result, nil
}

func checkUrl(url string) (bool, error) {
	fmt.Println("checking...")
	fmt.Println(url)
	cleanUrl := strings.Split(url, ":")[0]
	if filter.TestString(cleanUrl) {
		return CheckUrlWithBrain(cleanUrl)
	} else {
		return false, nil
	}
}

func IsMalwareRequestHttp() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTP")
		res, err := checkUrl(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

func IsMalwareRequestHttps() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTPS")
		res, err := checkUrl(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

var IsConnectToMalware goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	fmt.Println("connecting...")
	return goproxy.MitmConnect, host
}
