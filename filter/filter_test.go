package filter

import (
	"testing"

	"github.com/Ringloop/pisec/mocks"
	"github.com/stretchr/testify/require"
)

func TestUrlInBrainServer(t *testing.T) {
	//given
	mockClient := mocks.MockBrainClient{}
	mockRepo := mocks.MockRedisRepo{}

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
	mockClient := mocks.MockBrainClient{}
	mockRepo := mocks.MockRedisRepo{}

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
	mockClient := mocks.MockBrainClient{}
	mockRepo := mocks.MockRedisRepo{}
	mockFilter := mocks.MockUrlFilter{}

	phishUrl := "allowedUrl.com"

	mockRepo.Allowed = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}

func TestUrlInCacheDenied(t *testing.T) {
	//given
	mockClient := mocks.MockBrainClient{}
	mockRepo := mocks.MockRedisRepo{}
	mockFilter := mocks.MockUrlFilter{}

	phishUrl := "deniedUrl.com"

	mockRepo.Denied = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.False(t, res)
}

func TestUrlInCacheFalsePositive(t *testing.T) {
	//given
	mockClient := mocks.MockBrainClient{}
	mockRepo := mocks.MockRedisRepo{}
	mockFilter := mocks.MockUrlFilter{}

	phishUrl := "falsePositiveUrl.com"

	mockRepo.FalsePositives = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInBloomFilterFunc = func(string) bool {
		return true
	}

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}
