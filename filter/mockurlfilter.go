package filter

type MockUrlFilter struct {
	CheckInFilterFunc func(url string) bool
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	CheckInFilterFunc func(url string) bool
)

// Do is the mock client's `Do` func
func (m *MockUrlFilter) CheckUrlInBloom(url string) bool {
	return CheckInFilterFunc(url)
}
