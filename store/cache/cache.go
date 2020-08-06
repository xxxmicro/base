package cache

import(
	"github.com/xxxmicro/base/store"
	"github.com/xxxmicro/base/store/memory"
)

type Cache interface {
	store.Store
}

type cache struct {
	m store.Store		// memory store
	b store.Store		// backing store
	options store.Options
}

func NewStore(store store.Store, opts ...store.Option) store.Store {
	return &cache{
		m: memory.NewStore(opts...),
		b: store,
	}
}

func (c *cache) init(opts ...store.Option) error {
	for _, o := range opts {
		o(&c.options)
	}
	return nil
}

func (c *cache) Init(opts ...store.Option) error {
	if err := c.init(opts...); err != nil {
		return err
	}

	if err := c.m.Init(opts...); err != nil {
		return err
	}
	return c.b.Init(opts...)
}

func (c *cache) Options() store.Options {
	return c.options
}

func (c *cache) Get(key string, opts ...store.ReadOption) (*store.Record, error) {
	rec, err := c.m.Get(key, opts...)
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	if rec != nil {
		return rec, nil
	}

	rec, err = c.b.Get(key, opts...)
	if err == nil {
		if err := c.m.Set(rec); err != nil {
			return nil, err
		}
	}
	return rec, err
}

func (c *cache) Set(r *store.Record, opts ...store.WriteOption) error {
	if err := c.m.Set(r, opts...); err != nil {
		return err
	}
	return c.b.Set(r, opts...)
}

func (c *cache) Delete(key string, opts ...store.DeleteOption) error {
	if err := c.m.Delete(key, opts...); err != nil {
		return err
	}
	return c.b.Delete(key, opts...)
}

func (c *cache) List(opts ...store.ListOption) ([]string, error) {
	keys, err := c.m.List(opts...)
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	if len(keys) > 0 {
		return keys, nil
	}

	keys, err = c.b.List(opts...)
	if err == nil {
		for _, key := range keys {
			rec, err := c.b.Get(key)
			if err != nil {
				return nil, err
			}
			if err := c.m.Set(rec); err != nil {
				return nil, err
			}
		}
	}
	return keys, err
}

/*
func (c *cache) Incr(key string) (*store.Record, error) {
	if err := c.m.Incr(key); err != nil {
		return err
	}
	return c.b.Incr(key)
}

func (c *cache) IncrBy(key string, value int64) (*store.Record, error) {
	if err := c.m.IncrBy(key, value); err != nil {
		return err
	}
	return c.b.IncrBy(key, value)
}

func (c *cache) Exists(key string, opts ...store.ReadOption) (bool, error) {
	var has bool
	if has, err := c.m.Exists(key, value); err != nil {
		return err
	}

	if has {
		return true, nil
	}

	return c.b.Exists(key, opts...)
}
*/

func (c *cache) Close() error {
	if err := c.m.Close(); err != nil {
		return err
	}
	return c.b.Close()
}

func (c *cache) String() string {
	return "cache"
}