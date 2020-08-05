package store

import(
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	DefaultStore Store = new(noopStore)
)

type Store interface {
	Init(...Option) error
	Options() Options
	Get(key string, opts ...ReadOption) (*Record, error)
	Set(r *Record, opts ...WriteOption) error
	Delete(key string, opts ...DeleteOption) error	
	List(opts ...ListOption)([]string, error)
	Close() error
	String() string
}

type Record struct {
	Key string `json:"key"`
	Value []byte `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
	Expiry time.Duration `json:"expiry,omitempty"`
}