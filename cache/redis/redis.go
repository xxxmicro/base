package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/xxxmicro/base/cache"
	"reflect"
	"strings"
	"time"
)

type RedisCache struct {
	options cache.Options
	r       *redis.Pool
}

func NewCache(opts ...cache.Option) cache.Cache {
	options := cache.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return &RedisCache{
		options: options,
	}
}

func (m *RedisCache) Init(opts ...cache.Option) error {
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

func (m *RedisCache) Options() cache.Options {
	return m.options
}

func (m *RedisCache) connect() (*redis.Pool, error) {
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
					_ = c.Close()
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

func (m *RedisCache) prefix(key string) string {
	if len(m.options.Prefix) > 0 {
		return fmt.Sprintf("%s:%s", m.options.Prefix, key)
	} else {
		return key
	}
}

func (m *RedisCache) String() string {
	return "redis"
}

func (m *RedisCache) Get(key string, resultPtr interface{}, opts ...cache.ReadOption) error {
	readOpts := cache.ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}

	key = m.prefix(key)

	c := m.r.Get()
	defer c.Close()

	data, err := redis.Bytes(c.Do("GET", key))
	if err == redis.ErrNil {
		return cache.ErrNil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(data, resultPtr)
}

func (m *RedisCache) BatchGet(keys []string, resultsPtr interface{}, opts ...cache.ReadOption) error {
	readOpts := cache.ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}

	realKeys := make([]interface{}, len(keys))
	for i, key := range keys {
		realKeys[i] = m.prefix(key)
	}

	c := m.r.Get()
	defer c.Close()

	var replies interface{}
	var err error
	if replies, err = c.Do("MGET", realKeys...); err != nil {
		return err
	}

	sliceType := reflect.TypeOf(resultsPtr).Elem()
	valuePtrType := sliceType.Elem()
	valueType := valuePtrType.Elem()
	slice := reflect.MakeSlice(sliceType, 0, 0)
	for _, v := range replies.([]interface{}) {
		var data []byte
		if data, err = redis.Bytes(v, nil); err == redis.ErrNil {
			slice = reflect.Append(slice, reflect.Zero(valuePtrType))
			continue
		} else if err != nil {
			return err
		}
		valuePtr := reflect.New(valueType)
		if err = json.Unmarshal(data, valuePtr.Interface()); err != nil {
			return err
		}
		slice = reflect.Append(slice, valuePtr)
	}
	reflect.ValueOf(resultsPtr).Elem().Set(slice)
	return nil
}

func (m *RedisCache) Set(key string, value interface{}, opts ...cache.WriteOption) error {
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

func (m *RedisCache) Delete(key string, opts ...cache.DeleteOption) error {
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

func (m *RedisCache) BatchDelete(keys []string, opts ...cache.DeleteOption) error {
	deleteOptions := cache.DeleteOptions{}
	for _, o := range opts {
		o(&deleteOptions)
	}

	realKeys := make([]interface{}, len(keys))
	for i, key := range keys {
		realKeys[i] = m.prefix(key)
	}

	c := m.r.Get()
	defer c.Close()

	_, err := c.Do("DEL", realKeys...)
	return err
}

func (m *RedisCache) Close() error {
	return m.r.Close()
}
