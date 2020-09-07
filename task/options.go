package task

import "context"

type Option func(o *Options)

type Options struct {
	Context context.Context
}
