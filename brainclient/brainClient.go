package brainclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bits-and-blooms/bloom/v3"
)

type BrainClient interface {
	CheckUrl(url string) (bool, error)
	DownloadBloomFilter()
}

type Client struct {
	brainAddress       string
	indicatorsEndpoint string
	detailsEndpoint    string
}

func NewClient(brainAddr string, indicatorsEndpoint string, detailsEndpoint string) *Client {
	return &Client{
		brainAddress:       brainAddr,
		detailsEndpoint:    brainAddr + detailsEndpoint,
		indicatorsEndpoint: brainAddr + indicatorsEndpoint,
	}
}

func (client *Client) DownloadBloomFilter() *bloom.BloomFilter {

	var filter *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)

	//download the bloom filter from server
	res, err := http.Get(client.indicatorsEndpoint)
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

func (client *Client) isUrlInBrainRepo(buf *bytes.Buffer) (bool, error) {
	res, err := http.Post(client.detailsEndpoint, "application/json", buf)
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

func (client *Client) CheckUrl(url string) (bool, error) {

	var checkUrlReq bytes.Buffer
	err := CreateCheckUrlReq(url, &checkUrlReq)
	if err != nil {
		return false, err
	}

	return client.isUrlInBrainRepo(&checkUrlReq)

}
