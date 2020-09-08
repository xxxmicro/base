package redis

import (
	"context"
	"github.com/xxxmicro/base/cache"
)

type addrsKey struct{}
type passwordKey struct{}

func WithAddrs(addrs ...string) cache.Option {
	return func(o *cache.Options) {
		o.Context = context.WithValue(o.Context, addrsKey{}, addrs)
	}
}

func WithPassword(password string) cache.Option {
	return func(o *cache.Options) {
		o.Context = context.WithValue(o.Context, passwordKey{}, password)
	}
}