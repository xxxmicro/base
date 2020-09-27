package xxljob

import (
	"context"
	xxl "github.com/xxl-job/go-client"
	"github.com/xxxmicro/base/task"
)

type XXLTaskExecutor struct {
	opts     task.Options
	executor xxl.Executor
}

func NewXxlTaskExecutor(opts ...task.Option) task.TaskExecutor {
	options := task.Options{}

	for _, o := range opts {
		o(&options)
	}

	return &XXLTaskExecutor{
		opts: options,
	}
}

func (e *XXLTaskExecutor) Init(opts ...task.Option) {
	for _, o := range opts {
		o(&e.opts)
	}

	addr := e.opts.Context.Value(addressKey{}).(string)
	registerKey := e.opts.Context.Value(registerKey{}).(string)

	executor := xxl.NewExecutor(
		xxl.ServerAddr(addr),
		//xxl.AccessToken(xxlConfig.AccessToken), //请求令牌(默认为空)
		xxl.RegistryKey(registerKey),
	)

	executor.Init()

	e.executor = executor
}

func (e *XXLTaskExecutor) RegisterTask(pattern string, t task.ITask) {
	e.executor.RegTask(pattern, func(c context.Context, param *xxl.RunReq) {
		req := &task.RunReq{
			JobID:                 param.JobID,
			ExecutorHandler:       param.ExecutorHandler,
			ExecutorParams:        param.ExecutorParams,
			ExecutorBlockStrategy: param.ExecutorBlockStrategy,
			ExecutorTimeout:       param.ExecutorTimeout,
			LogID:                 param.LogID,
			LogDateTime:           param.LogDateTime,
			GlueType:              param.GlueType,
			GlueSource:            param.GlueSource,
			GlueUpdatetime:        param.GlueUpdatetime,
			BroadcastIndex:        param.BroadcastIndex,
			BroadcastTotal:        param.BroadcastTotal,
		}

		t.Execute(c, req)
	})
}

func (e *XXLTaskExecutor) Start() {
	go e.executor.Run()
}
