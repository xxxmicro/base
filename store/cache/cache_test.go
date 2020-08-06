package cache

import(
	"testing"
	"log"
	"time"
	"github.com/xxxmicro/base/store"
	"github.com/xxxmicro/base/store/redis"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cf := NewStore(redis.NewStore())
	cf.Init()
	cfInt := cf.(*cache)

	record := &store.Record{
		Key: "key1",
		Value: []byte("foo"),
		Expiry: time.Millisecond * 1000,
	}

	err := cf.Delete(record.Key)
	assert.NoError(t, err)

	r, err := cf.Get("key1")
	assert.Error(t, err, "Unexpected record")
	cfInt.b.Set(record)

	time.Sleep(time.Millisecond * 500)
	r, err = cf.Get("key1")
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 2000)
	r, err = cf.Get(record.Key)
	assert.Error(t, err, "Expected no records in redis store")

	log.Print(r)
}