package filter

import (
	"fmt"
	"strings"

	"github.com/Ringloop/pisec/brainclient"
	"github.com/Ringloop/pisec/cache"
	"github.com/bits-and-blooms/bloom/v3"
)

type PisecUrlFilter struct {
	Client      brainclient.BrainClient
	Repo        *cache.RedisRepository
	bloomFilter *bloom.BloomFilter
}

func NewPisecUrlFilter(c brainclient.BrainClient, r *cache.RedisRepository) *PisecUrlFilter {
	f := &PisecUrlFilter{
		Client:      c,
		Repo:        r,
		bloomFilter: c.DownloadBloomFilter(),
	}
	return f
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

	if allow, err := psFilter.Repo.IsAllow(cleanUrl); err == nil {
		if allow {
			return true, nil //URL is allowed
		}
	} else { //err != nil
		return false, err
	}

	if falsePositive, err := psFilter.Repo.IsFalsePositive(cleanUrl); err == nil {
		if falsePositive {
			return true, nil //URL is a well known FALSE POSITIVE
		}
	} else { //err != nil
		return false, err
	}

	if deny, err := psFilter.Repo.IsDeny(cleanUrl); err == nil {
		if deny {
			return false, nil //URL is a well known POSITIVE
		}
	} else { //err != nil
		return false, err
	}

	return psFilter.Client.CheckUrl(url)

}
