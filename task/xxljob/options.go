package xxljob

import (
	"context"
	"github.com/xxxmicro/base/task"
)

type addressKey struct{}
type registerKey struct{}

func WithAddr(addr string) task.Option {
	return func(o *task.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, addressKey{}, addr)
	}
}

func WithRegisterKey(register string) task.Option {
	return func(o *task.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, registerKey{}, register)
	}
}
