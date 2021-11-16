package redis_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xxxmicro/base/cache"
	"github.com/xxxmicro/base/cache/redis"
	"log"
	"testing"
	"time"
)

type User struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func TestBasic(t *testing.T) {
	s := redis.NewCache(redis.WithAddrs(":6379"))
	s.Init()

	key := "hello"
	value := &User{Name: "李小龙", CreatedAt: time.Now()}

	err := s.Set(key, value, cache.WriteExpiry(time.Millisecond*1000))
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)

	value1 := &User{}
	err = s.Get(key, value1)
	assert.NoError(t, err)
	log.Print(value1)

	time.Sleep(time.Millisecond * 2000)

	value2 := &User{}
	err = s.Get(key, value2)
	assert.Equal(t, cache.ErrNil, err)
}

func TestBatch(t *testing.T) {
	s := redis.NewCache(redis.WithAddrs(":6379"))
	s.Init()

	var err error

	key1 := "bruceli"
	value1 := &User{Name: "李小龙", CreatedAt: time.Now()}

	key2 := "jackychen"
	value2 := &User{Name: "成龙", CreatedAt: time.Now()}

	err = s.Set(key1, value1)
	assert.NoError(t, err)

	err = s.Set(key2, value2)
	assert.NoError(t, err)

	users := make([]*User, 0)
	err = s.BatchGet([]string{key1, "fatchow", key2}, &users)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(users))
	assert.Equal(t, value1.Name, users[0].Name)
	assert.Nil(t, users[1])
	assert.Equal(t, value2.Name, users[2].Name)
	t.Log(users[0])
	t.Log(users[1])
	t.Log(users[2])

	err = s.BatchDelete([]string{key1, "fatchow", key2})
	assert.NoError(t, err)

	gotValue := &User{}
	err = s.Get(key1, gotValue)
	assert.Equal(t, cache.ErrNil, err)

	err = s.Get(key2, gotValue)
	assert.Equal(t, cache.ErrNil, err)
}
