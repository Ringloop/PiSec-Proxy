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

func TestUrlInBloomFilter(t *testing.T) {
	//given
	repo := cache.NewRedisClient()
	err = repo.InitRepository()
	require.Nil(t, err)

	phishUrl := "newPhishingUrl.com"

	var bloomFilter *bloom.BloomFilter = bloom.NewWithEstimates(1000000, 0.01)
	bloomFilter.AddString(phishUrl)

	filter := NewPisecUrlFilter(repo, bloomFilter, indicatorsEndpoint)

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.False(t, res)
}
