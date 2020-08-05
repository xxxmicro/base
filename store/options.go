package store

import(
	"context"
	"time"
)

type Options struct {
	Addrs 					[]string
	Database				string
	Table						string
	Password 				string
	Context 				context.Context
}

type Option func(o *Options)

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Database allows multiple isolated stores to be kept in one backend, if supported.
func Database(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

// Table is analagous to a table in database backends or a key prefix in KV backends
func Table(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithContext(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

type ReadOptions struct {
	Database, Table string
	Prefix bool
	Suffix bool
	Limit uint
	Offset uint
}

type ReadOption func(r *ReadOptions)

func ReadFrom(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	} 
}

func ReadPrefix() ReadOption {
	return func(r *ReadOptions) {
		r.Prefix = true
	}
}

func ReadSuffix() ReadOption {
	return func(r *ReadOptions) {
		r.Suffix = true
	}
}

func ReadLimit(l uint) ReadOption {
	return func(r *ReadOptions) {
		r.Limit = l
	}
}

func ReadOffset(o uint) ReadOption {
	return func(r *ReadOptions) {
		r.Offset = o
	}
}

type WriteOptions struct {
	Database, Table string
	Expiry time.Time
	TTL time.Duration
}

type WriteOption func(w *WriteOptions)

func WriteTo(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

func WriteExpiry(t time.Time) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
	}
}

func WriteTTL(d time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.TTL = d
	}
}

type DeleteOptions struct {
	Database, Table string
}

type DeleteOption func(d *DeleteOptions)

func DeleteFrom(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

type ListOptions struct {
	Database, Table string
	Prefix string
	Suffix string
	Limit uint
	Offset uint
}

type ListOption func(l *ListOptions)

func ListFrom(database, table string) ListOption {
	return func(lo *ListOptions) {
		lo.Database = database
		lo.Table = table
	}
}

func ListPrefix(p string) ListOption {
	return func(lo *ListOptions) {
		lo.Prefix = p
	}
}

func ListSuffix(s string) ListOption {
	return func(lo *ListOptions) {
		lo.Suffix = s
	}
}

func ListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

func ListOffset(o uint) ListOption {
	return func(lo *ListOptions) {
		lo.Offset = o
	}
}