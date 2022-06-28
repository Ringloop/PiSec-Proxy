package filter

type MockRedisRepo struct {
	Allowed        map[string]string
	Denied         map[string]string
	FalsePositives map[string]string
}

func (r MockRedisRepo) IsAllow(url string) (bool, error) {
	_, isAllowed := r.Allowed[url]
	return isAllowed, nil
}

func (r MockRedisRepo) IsDeny(url string) (bool, error) {
	_, isDenied := r.Denied[url]
	return isDenied, nil
}

func (r MockRedisRepo) IsFalsePositive(url string) (bool, error) {
	_, isFalsePositive := r.FalsePositives[url]
	return isFalsePositive, nil
}
