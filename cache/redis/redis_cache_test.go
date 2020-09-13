package redis_test

import(
	"testing"
	"time"
	"log"
	"github.com/xxxmicro/base/cache"
	"github.com/xxxmicro/base/cache/redis"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func TestBasic(t *testing.T) {
	s := redis.NewCache(redis.WithAddrs(":6379"))
	s.Init()

	key := "hello"
	value := &User{Name:"李小龙", CreatedAt: time.Now()}

	err := s.Set(key, value, cache.WriteExpiry(time.Millisecond * 1000))
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)

	value1 := &User{}
	err = s.Get(key, value1)
	assert.NoError(t, err)
	log.Print(value1)

	time.Sleep(time.Millisecond * 2000)

	value2 := &User{}
	err = s.Get(key, value2)
	assert.Error(t, err, "Expected no records in redis store")
}