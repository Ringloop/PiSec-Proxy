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

	"github.com/Ringloop/pisec/cache"
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

var repo *cache.RedisRepository

func main() {

	//setup the REDIS cache
	repo = cache.NewRedisClient()

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

	err = filter.UnmarshalJSON(jsonRes)
	if err != nil {
		panic(err)
	}

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

func createCheckUrlReq(url string, buf *bytes.Buffer) error {
	rq := CheckUrlRequest{
		Url: url,
	}
	err := json.NewEncoder(buf).Encode(rq)
	return err
}

func isUrlInBrainRepo(buf *bytes.Buffer) (bool, error) {
	endpoint := serverAddress + detailsEndpoint
	res, err := http.Post(endpoint, "application/json", buf)
	if err != nil {
		return false, err
	}

	jsonRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var checkUrlRes CheckUrlResponse
	err = json.Unmarshal(jsonRes, &checkUrlRes)

	if err != nil {
		return false, err
	}

	return checkUrlRes.Result, nil

}

func CheckUrlWithBrain(url string) (bool, error) {

	var checkUrlReq bytes.Buffer
	err := createCheckUrlReq(url, &checkUrlReq)
	if err != nil {
		return false, err
	}

	isUrlInRepo, err := isUrlInBrainRepo(&checkUrlReq)
	if err != nil {
		return false, err
	}

	if isUrlInRepo {
		repo.AddDeny(url)
		return true, nil
	} else {
		repo.AddFalsePositive(url)
		return false, nil
	}

}

/*
This function says if the navigation to the passed URL is allowed or not.
Cases are as following (order is important)
  - Url is NOT found in the Bloom Filter: return TRUE because the URL is not in the repository, for sure.
  	All the other cases requires that the URL has been found in the Bloom Filter
  - URL is in ALLOW cache: return TRUE because the URL is a malicious one, but the user has already allowed the navigation through this
  - URL is in FALSE cache: return TRUE because the URL is a false positive of the Bloom Filter, already checked
  - URL is in DENY cache: return FALSE because the URL is a malicious one, and it has been already checked with server and blocked
  - Outcome is dubious, so we need to check this result with Brain server, cache will be updated accordingly
*/
func shallYouPass(url string) (bool, error) {
	fmt.Println("checking...")
	fmt.Println(url)
	cleanUrl := strings.Split(url, ":")[0]

	if !filter.TestString(cleanUrl) {
		return true, nil //URL is NOT present, for sure
	}

	if allow, err := repo.IsAllow(cleanUrl); err == nil {
		if allow {
			return true, nil //URL is allowed
		}
	} else { //err != nil
		return false, err
	}

	if falsePositive, err := repo.IsFalsePositive(cleanUrl); err == nil {
		if falsePositive {
			return true, nil //URL is a well known FALSE POSITIVE
		}
	} else { //err != nil
		return false, err
	}

	if deny, err := repo.IsDeny(cleanUrl); err == nil {
		if deny {
			return false, nil //URL is a well known POSITIVE
		}
	} else { //err != nil
		return false, err
	}
	return CheckUrlWithBrain(cleanUrl)

}

func IsMalwareRequestHttp() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTP")
		res, err := shallYouPass(strings.Split(req.Host, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
}

func IsMalwareRequestHttps() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println("Inside HTTPS")
		res, err := shallYouPass(strings.Split(req.Host, ":")[0])
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
