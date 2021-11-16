package cache

import "errors"

type Cache interface {
	Init(opts ...Option) error
	Options() Options
	Get(key string, resultPtr interface{}, opts ...ReadOption) error
	BatchGet(keys []string, resultsPtr interface{}, opts ...ReadOption) error
	Set(key string, value interface{}, opts ...WriteOption) error
	Delete(key string, opts ...DeleteOption) error
	BatchDelete(keys []string, opts ...DeleteOption) error
}

var ErrNil = errors.New("nil returned")
