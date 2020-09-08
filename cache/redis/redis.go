package redis

import(
	"context"
	"fmt"
	"time"
	"strings"
	"github.com/xxxmicro/base/cache"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

type redisCache struct {
	options cache.Options
	r *redis.Pool
}

func NewCache(opts ...cache.Option) cache.Cache {
	options := cache.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &redisCache{
		options: options,
	}
}

func (m *redisCache) Init(opts ...cache.Option) error {
	for _, o := range opts {
		o(&m.options)
	}

	r, err := m.connect()
	if err != nil {
		return err
	}
	m.r = r
	return nil
}

func (m *redisCache) Options() cache.Options {
	return m.options
}

func (m *redisCache) connect() (*redis.Pool, error) {
	addrs := m.options.Context.Value(addrsKey{}).([]string)
	if len(addrs) == 0 {
		addrs = []string{":6379"}
	}

	password, _ := m.options.Context.Value(passwordKey{}).(string)

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", strings.Join(addrs, ","))
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return pool, nil
}

func (m *redisCache) prefix(key string) string {
	return fmt.Sprintf("%s:%s", m.options.Prefix, key)
}

func (m *redisCache) String() string {
	return "redis"
}

func (m *redisCache) Get(key string, resultPtr interface{}, opts ...cache.ReadOption) error {
	readOpts := cache.ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}

	key = m.prefix(key)

	c := m.r.Get()
	defer c.Close()
	data, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, resultPtr)
}

func (m *redisCache) Set(key string, value interface{}, opts ...cache.WriteOption) error {
	writeOpts := cache.WriteOptions{}
	for _, o := range opts {
		o(&writeOpts)
	}

	key = m.prefix(key)

	c := m.r.Get()
	defer c.Close()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if writeOpts.Expiry > 0 {
		_, err = c.Do("SETEX", key, writeOpts.Expiry.Seconds(), data)
	} else {
		_, err = c.Do("SET", key, data)
	}
	return err
}

func (m *redisCache) Close() error {
	return m.r.Close()
}

func (m *redisCache) Delete(key string, opts ...cache.DeleteOption) error {
	deleteOptions := cache.DeleteOptions{}
	for _, o := range opts {
		o(&deleteOptions)
	}

	key = m.prefix(key)

	c := m.r.Get()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}