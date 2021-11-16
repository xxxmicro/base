package cache

import (
	"context"
	"time"
)

type Options struct {
	Context context.Context
	Prefix  string
}

type Option func(o *Options)

type ReadOptions struct {
}

type ReadOption func(o *ReadOptions)

type WriteOptions struct {
	Expiry time.Duration
}

type WriteOption func(o *WriteOptions)

func WriteExpiry(t time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
	}
}

type DeleteOptions struct {
}

type DeleteOption func(o *DeleteOptions)
