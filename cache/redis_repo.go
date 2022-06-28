package cache

import (
	"time"

	"github.com/go-redis/redis"
)

const DAYS = 24 * time.Hour
const EXPIRE_DAYS = 7 * int(DAYS)

const PISEC_KEY = "PISEC:"
const ALLOW_KEY = PISEC_KEY + "ALLOW:"
const DENY_KEY = PISEC_KEY + "DENY:"
const FALSE_POSITIVE_KEY = PISEC_KEY + "FALSE_POSITIVE:"

const DEFAULT_PORT = "6379"
const TEST_PORT = "6378"

type RepoClient interface {
	IsAllow(url string) (bool, error)
	IsDeny(url string) (bool, error)
	IsFalsePositive(url string) (bool, error)
}

type RedisRepository struct {
	client       *redis.Client
	dataDuration int
}

func NewTestClient() *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + DEFAULT_PORT,
		Password: "test", // no password set
		DB:       0,      // use default DB
	})
	return &RedisRepository{
		client:       rdb,
		dataDuration: 0,
	}
}

func NewRedisClient() *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + DEFAULT_PORT,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // no password set
		DB:       0,                                  // use default DB
	})

	return &RedisRepository{
		client:       rdb,
		dataDuration: EXPIRE_DAYS}
}

func (repo *RedisRepository) AddAllow(url string) error {
	return repo.client.Set(ALLOW_KEY+url, url, time.Duration(repo.dataDuration)).Err()
}

func (repo *RedisRepository) AddDeny(url string) error {
	return repo.client.Set(DENY_KEY+url, url, time.Duration(repo.dataDuration)).Err()
}

func (repo *RedisRepository) AddFalsePositive(url string) error {
	return repo.client.Set(FALSE_POSITIVE_KEY+url, url, time.Duration(repo.dataDuration)).Err()
}

func (repo *RedisRepository) IsAllow(url string) (bool, error) {
	res, err := repo.client.Get(ALLOW_KEY + url).Result()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	return res == url, nil

}

func (repo *RedisRepository) IsDeny(url string) (bool, error) {
	res, err := repo.client.Get(DENY_KEY + url).Result()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	return res == url, nil

}

func (repo *RedisRepository) IsFalsePositive(url string) (bool, error) {
	res, err := repo.client.Get(FALSE_POSITIVE_KEY + url).Result()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	return res == url, nil

}

func (repo *RedisRepository) InitRepository() error {
	return repo.client.FlushDB().Err()
}

func (repo *RedisRepository) GetRepoSize() (int, error) {
	res := repo.client.DbSize()
	return int(res.Val()), res.Err()
}
