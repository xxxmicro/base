package redis

import(
	"testing"
	"time"
	"log"
	"github.com/xxxmicro/base/store"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	s := NewStore(store.Addrs(":6379"))
	s.Init()

	record := &store.Record{
		Key: "hello",
		Value: []byte("world"),
		Expiry: time.Millisecond * 1000,
	}

	err := s.Set(record)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)
	r, err := s.Get(record.Key)
	assert.NoError(t, err)
	log.Print(r)

	time.Sleep(time.Millisecond * 2000)

	r, err = s.Get(record.Key)
	assert.Error(t, err, "Expected no records in redis store")
}