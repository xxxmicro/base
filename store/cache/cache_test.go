package cache

import(
	"os"
	"path/filepath"
	"testing"
	"github.com/xxxmicro/base/store"
)

func cleanup(db string, s store.Store) {
	s.Close()
	dir := filepath.Join(file.DefaultDir, db + "/")
	os.RemoveAll(dir)
}

func TestGet(t *testing.T) {
	cf := NewStore(redis.NewStore())
	cf.Init()
	cfInt := cf.(*cache)

	_, err := cf.Get("key1")
	assert.Error(t, err, "Unexpected record")
	cfInt.b.Write(&store.Record{
		Key: "key1",
		Value: []byte("foo"),
	})
	rec, err := cf.Get("key1")
	assert.NoError(t, err)
}