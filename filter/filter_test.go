package filter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var err error
var indicatorsEndpoint = os.Getenv("PISEC_BRAIN_ADDR") + "/api/v1/indicators/details"

func TestUrlInBrainServer(t *testing.T) {
	//given
	mockClient := MockBrainClient{}
	mockRepo := MockRedisRepo{}

	phishUrl := "evil.com"

	filter := NewPisecUrlFilter(mockClient, mockRepo)

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}

func TestUrlNotInBloomFilter(t *testing.T) {
	//given
	mockClient := MockBrainClient{}
	mockRepo := MockRedisRepo{}

	filter := NewPisecUrlFilter(mockClient, mockRepo)

	phishUrl := "bloomfiltertest.com"

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}

func TestUrlInCacheAllowed(t *testing.T) {
	//given
	mockClient := MockBrainClient{}
	mockRepo := MockRedisRepo{}

	filter := NewPisecUrlFilter(mockClient, mockRepo)

	phishUrl := "allowedUrl.com"

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}
