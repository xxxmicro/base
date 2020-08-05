package redis

import(
	"time"
	"strings"
	"path/filepath"
	"errors"
	"github.com/xxxmicro/base/store"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

type storeRecord struct {
	Key string											`json:"key"`
	Value []byte										`json:'value'`
	Metadata map[string]interface{}	`json:"metadata"`
	ExpiresAt int64									`json:"expires_at"`
}

type redisStore struct {
	options store.Options
	r *redis.Pool
}

func NewStore(opts ...store.Option) store.Store {
	s := &redisStore{
		options: store.Options{},
	}

	for _, o := range opts {
		o(&s.options)
	}

	return s
}

func (m *redisStore) Init(opts ...store.Option) error {
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

func (m *redisStore) Options() store.Options {
	return m.options
}

func (m *redisStore) prefix(database, table string) string {
	if len(database) == 0 {
		database = m.options.Database
	}

	if len(table) == 0 {
		table = m.options.Table
	}
	return filepath.Join(database, table)
}

func (m *redisStore) connect() (*redis.Pool, error) {
	if len(m.options.Addrs) == 0 {
		m.options.Addrs = []string{":6379"}
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
						c, err := redis.Dial("tcp", strings.Join(m.options.Addrs, ","))
						if err != nil {
										return nil, err
						}

						if m.options.Password != "" {
										if _, err := c.Do("AUTH", m.options.Password); err != nil {
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

func (m *redisStore) String() string {
	return "redis"
}

func (m *redisStore) Get(key string, opts ...store.ReadOption) (*store.Record, error) {
	readOpts := store.ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}

	prefix := m.prefix(readOpts.Database, readOpts.Table)
	key = filepath.Join(prefix, key)

	c := m.r.Get()
	defer c.Close()
	data, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		return nil, err
	}

	/*
	if len(data) == 0 {
		return nil, nil
	}
	*/
	
	var r *storeRecord
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	newRecord := &store.Record{}
	newRecord.Key = strings.TrimPrefix(r.Key, prefix+"/")
	newRecord.Value = make([]byte, len(r.Value))
	newRecord.Metadata = make(map[string]interface{})

	// copy the value into the new record
	copy(newRecord.Value, r.Value)

	// check if we need to set the expiry
	if r.ExpiresAt != 0 {
		newRecord.Expiry = time.Until(time.Unix(r.ExpiresAt, 0))
	}

	// copy in the metadata
	for k, v := range r.Metadata {
		newRecord.Metadata[k] = v
	}

	return newRecord, err
}

func (m *redisStore) Set(record *store.Record, opts ...store.WriteOption) error {
	writeOpts := store.WriteOptions{}
	for _, o := range opts {
		o(&writeOpts)
	}

	prefix := m.prefix(writeOpts.Database, writeOpts.Table)
	key := filepath.Join(prefix, record.Key)

	if !writeOpts.Expiry.IsZero() {
		record.Expiry = time.Until(writeOpts.Expiry)
	}

	if writeOpts.TTL != 0 {
		record.Expiry = writeOpts.TTL
	}

	// copy the incoming record and then convert the expiry in to a hard timestamp
	i := storeRecord{}
	i.Key = record.Key
	i.Value = make([]byte, len(record.Value))
	i.Metadata = make(map[string]interface{})

	// copy the value
	copy(i.Value, record.Value)

	// set the expiry
	if record.Expiry != 0 {
		i.ExpiresAt = time.Now().Add(record.Expiry).Unix()
	}

	// set the metadata
	for k, v := range record.Metadata {
		i.Metadata[k] = v
	}

	c := m.r.Get()
	defer c.Close()

	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if record.Expiry > 0 {
		_, err = c.Do("SETEX", key, record.Expiry.Seconds(), data)
	} else {
		_, err = c.Do("SET", key, data)
	}
	return err
}

func (m *redisStore) Close() error {
	return m.r.Close()
}

func (m *redisStore) Delete(key string, opts ...store.DeleteOption) error {
	deleteOptions := store.DeleteOptions{}
	for _, o := range opts {
		o(&deleteOptions)
	}

	prefix := m.prefix(deleteOptions.Database, deleteOptions.Table)
	key = filepath.Join(prefix, key)

	c := m.r.Get()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}

func (m *redisStore) List(opts ...store.ListOption) ([]string, error) {
	listOptions := store.ListOptions{}
	for _, o := range opts {
		o(&listOptions)
	}

	// prefix := m.prefix(listOptions.Database, listOptions.Table)
	// TODO	
	// keys := make([]string, 0)

	return nil, errors.New("not impl yet")
}