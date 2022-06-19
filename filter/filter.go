package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ringloop/pisec/cache"
	"github.com/bits-and-blooms/bloom/v3"
)

type PisecUrlFilter struct {
	brainEndpoint   string
	repo            *cache.RedisRepository
	bloomFilter     *bloom.BloomFilter
	detailsEndpoint string
}

func NewPisecUrlFilter(repo *cache.RedisRepository, bloomFilter *bloom.BloomFilter, detailsEndpoint string) *PisecUrlFilter {

	return &PisecUrlFilter{
		repo:            repo,
		bloomFilter:     bloomFilter,
		detailsEndpoint: detailsEndpoint,
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
func (psFilter *PisecUrlFilter) ShallYouPass(url string) (bool, error) {
	fmt.Println("checking...")
	fmt.Println(url)
	cleanUrl := strings.Split(url, ":")[0]

	if !psFilter.bloomFilter.TestString(cleanUrl) {
		return true, nil //URL is NOT present, for sure
	}

	if allow, err := psFilter.repo.IsAllow(cleanUrl); err == nil {
		if allow {
			return true, nil //URL is allowed
		}
	} else { //err != nil
		return false, err
	}

	if falsePositive, err := psFilter.repo.IsFalsePositive(cleanUrl); err == nil {
		if falsePositive {
			return true, nil //URL is a well known FALSE POSITIVE
		}
	} else { //err != nil
		return false, err
	}

	if deny, err := psFilter.repo.IsDeny(cleanUrl); err == nil {
		if deny {
			return false, nil //URL is a well known POSITIVE
		}
	} else { //err != nil
		return false, err
	}
	return psFilter.checkUrlWithBrain(cleanUrl)

}

func (psFilter *PisecUrlFilter) isUrlInBrainRepo(buf *bytes.Buffer) (bool, error) {
	res, err := http.Post(psFilter.detailsEndpoint, "application/json", buf)
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

func (psFilter *PisecUrlFilter) checkUrlWithBrain(url string) (bool, error) {

	var checkUrlReq bytes.Buffer
	err := CreateCheckUrlReq(url, &checkUrlReq)
	if err != nil {
		return false, err
	}

	isUrlInRepo, err := psFilter.isUrlInBrainRepo(&checkUrlReq)
	if err != nil {
		return false, err
	}

	if isUrlInRepo {
		psFilter.repo.AddDeny(url)
		return true, nil
	} else {
		psFilter.repo.AddFalsePositive(url)
		return false, nil
	}

}
