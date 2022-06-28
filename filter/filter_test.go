package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
	mockFilter := MockUrlFilter{}

	phishUrl := "allowedUrl.com"

	mockRepo.Allowed = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInFilterFunc = func(string) bool {
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
	mockClient := MockBrainClient{}
	mockRepo := MockRedisRepo{}
	mockFilter := MockUrlFilter{}

	phishUrl := "deniedUrl.com"

	mockRepo.Denied = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInFilterFunc = func(string) bool {
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
	mockClient := MockBrainClient{}
	mockRepo := MockRedisRepo{}
	mockFilter := MockUrlFilter{}

	phishUrl := "falsePositiveUrl.com"

	mockRepo.FalsePositives = map[string]string{
		phishUrl: "true",
	}

	mockFilter.CheckInFilterFunc = func(string) bool {
		return true
	}

	filter := NewPisecUrlFilter(mockClient, mockRepo)
	filter.CheckInFilterFunc = func(string) bool {
		return true
	}

	//when
	res, err := filter.ShallYouPass(phishUrl)

	//then
	require.Nil(t, err)
	require.True(t, res)
}
