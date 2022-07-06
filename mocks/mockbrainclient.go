package mocks

import "github.com/bits-and-blooms/bloom/v3"

type MockBrainClient struct {
	ListedUrls map[string]struct{}
}

func (f MockBrainClient) CheckUrl(url string) (bool, error) {
	_, isUrlListed := f.ListedUrls[url]
	return isUrlListed, nil
}

func (f MockBrainClient) DownloadBloomFilter() *bloom.BloomFilter {
	var bf *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)
	bf.AddString("evil.com")
	return bf
}
