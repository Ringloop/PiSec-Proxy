package mocks

import (
	"github.com/Ringloop/pisec/brainclient"
	"github.com/Ringloop/pisec/cache"
)

type MockUrlFilter struct {
	CheckInFilterFunc func(url string) bool
	Client            brainclient.BrainClient
	Repo              cache.RepoClient
}

var (
	// CheckInFilterFunc fetches the mock client's `CheckUrlInBloom` func
	CheckInFilterFunc func(url string) bool
)

// CheckUrlInBloom is the mock client's `CheckUrlInBloom` func
func (m *MockUrlFilter) CheckUrlInBloom(url string) bool {
	return CheckInFilterFunc(url)
}
