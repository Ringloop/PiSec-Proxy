package mocks

import "github.com/bits-and-blooms/bloom/v3"

type MockBrainClient struct {
}

func (f MockBrainClient) CheckUrl(url string) (bool, error) {
	if url == "evil.com" {
		return true, nil
	}

	return false, nil
}

func (f MockBrainClient) DownloadBloomFilter() *bloom.BloomFilter {
	var bf *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)
	bf.AddString("evil.com")
	return bf
}
