package mocks

type MockUrlFilter struct {
	CheckInFilterFunc func(url string) bool
}

// CheckUrlInBloom is the mock client's `CheckUrlInBloom` func
func (m *MockUrlFilter) CheckUrlInBloom(url string) bool {
	return m.CheckInFilterFunc(url)
}
