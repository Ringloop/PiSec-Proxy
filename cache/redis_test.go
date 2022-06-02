package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisClient(t *testing.T) {
	//given
	repo := NewRedisClient()

	//when
	pong, err := repo.client.Ping().Result()

	//then
	require.Nil(t, err)
	require.Equal(t, pong, "PONG")
}

func TestRedisInit(t *testing.T) {

	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	err := repo.AddAllow("ALLOW_EXAMPLE")
	require.Nil(t, err)

	//then
	val, err := repo.GetRepoSize()
	require.Nil(t, err)
	require.Equal(t, val, 1)

	err = repo.InitRepository()
	require.Nil(t, err)

	val, err = repo.GetRepoSize()
	require.Nil(t, err)
	require.Equal(t, val, 0)

}

func TestRedisInsertAllow(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	err := repo.AddAllow("ALLOW_EXAMPLE")
	require.Nil(t, err)

	res, err := repo.IsAllow("ALLOW_EXAMPLE")
	require.Nil(t, err)

	//then
	require.True(t, res)
}

func TestRedisInsertDeny(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	err := repo.AddDeny("DENY_EXAMPLE")
	require.Nil(t, err)

	res, err := repo.IsDeny("DENY_EXAMPLE")
	//then
	require.Nil(t, err)
	require.True(t, res)
}

func TestRedisInsertFalsePositive(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	err := repo.AddFalsePositive("FALSE_POSITIVE_EXAMPLE")
	require.Nil(t, err)

	res, err := repo.IsFalsePositive("FALSE_POSITIVE_EXAMPLE")
	//then
	require.Nil(t, err)
	require.True(t, res)
}

func TestRedisMissingDeny(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	res, err := repo.IsDeny("DENY_EXAMPLE")

	//then
	require.Nil(t, err)
	require.False(t, res)
}

func TestRedisMissingAllow(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	res, err := repo.IsDeny("ALLOW_EXAMPLE")

	//then
	require.Nil(t, err) //Object not found is a special error case
	require.False(t, res)
}

func TestRedisMissingFalsePositive(t *testing.T) {
	//given
	repo := NewRedisClient()
	repo.InitRepository()

	//when
	res, err := repo.IsFalsePositive("FALSE_POSITIVE_EXAMPLE")

	//then
	require.Nil(t, err) //Object not found is a special error case
	require.False(t, res)
}
