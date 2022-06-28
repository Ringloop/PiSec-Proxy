package filter

import (
	"os"
	"testing"

	"github.com/Ringloop/pisec/cache"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/require"
)

var err error
var indicatorsEndpoint = os.Getenv("PISEC_BRAIN_ADDR") + "/api/v1/indicators/details"

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

func TestUrlInBloomFilter(t *testing.T) {
	//given
	mockClient := MockBrainClient{}

	repo := cache.NewRedisClient()
	err = repo.InitRepository()
	require.Nil(t, err)

	phishUrl := "evil.com"

	filter := NewPisecUrlFilter(mockClient, repo)

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}
