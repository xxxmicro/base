package memory

import(
	"os"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/xxxmicro/base/store"
)

func TestMemoryReInit(t *testing.T) {
	s := NewStore(store.Table("aaa"))
	s.Init(store.Table(""))
	if len(s.Options().Table) > 0 {
		t.Error("Init didn't reinitialize the store")
	}
}

func TestMemoryBasic(t *testing.T) {
	s := NewStore()
	s.Init()
	basictest(s, t)
}

func TestMemoryPrefix(t *testing.T) {
	s := NewStore()
	s.Init(store.Table("some-prefix"))
	basictest(s, t)
}

func TestMemoryNamespace(t *testing.T) {
	s := NewStore()
	s.Init(store.Database("some-namespace"))
	basictest(s, t)
}

func TestMemoryNamespacePrefix(t *testing.T) {
	s := NewStore()
	s.Init(store.Table("some-prefix"), store.Database("some-namespace"))
	basictest(s, t)
}

func basictest(s store.Store, t *testing.T) {
	if len(os.Getenv("IN_TRAVIS_CI")) == 0 {
		t.Logf("Testing store %s, with options %# v\n", s.String(), pretty.Formatter(s.Options()))
	}

	// Read and Write an expiring record
	if err := s.Set(&store.Record{
		Key: "Hello",
		Value: []byte("World"),
		Expiry: time.Millisecond * 100,
	}); err != nil {
		t.Error(err)
	}

	if r, err := s.Get("Hello"); err != nil {
		t.Error(err)
	} else {
		if r.Key != "Hello" {
			t.Errorf("Expected %s, got %s", "Hello", r.Key)
		}

		if string(r.Value) != "World" {
			t.Errorf("Expected %s, got %s", "World", r.Value)	
		}
	}

	time.Sleep(time.Millisecond * 200)

	if _, err := s.Get("Hello"); err != store.ErrNotFound {
		t.Errorf("Expected %# v, got %# v", store.ErrNotFound, err)
	}

	time.Sleep(time.Millisecond * 100)
	if results, err := s.List(); err != nil {
		t.Errorf("List failed: %s", err)
	} else {
		if len(results) != 0 {
			t.Error("Expiry options were not effective")
		}
	}
	
	s.Close() // reset the store
}