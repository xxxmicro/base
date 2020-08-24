package jaeger

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	xxxmicro_opentracing "github.com/xxxmicro/base/opentracing"
)

func NewTracerProvider(c *cli.Context, config config.Config) (tracer opentracing.Tracer, err error) {
	serviceName := config.Get("service", "name").String("")
	if len(serviceName) == 0 {
		serviceName = c.String("server_name")
	}

	if len(serviceName) == 0 {
		serviceName = "unamed"
	}

	agentAddr := config.Get("jaeger", "agent", "addr").String("localhost:6831")

	metricsFactory := prometheus.New()

	// 根据配置初始化Tracer 返回Closer
	tracer, _, err = (&jconfig.Configuration{
		ServiceName: serviceName,
		Disabled:    false,
		Sampler: &jconfig.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			// param的值在0到1之间，设置为1则将所有的Operation输出到Reporter
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentAddr,
		},
	}).NewTracer(jconfig.Metrics(metricsFactory))

	if err != nil {
		err = errors.Wrap(err, "create jaeger tracer error")
		logger.Error(err)
	}

	// 设置全局Tracer - 如果不设置将会导致上下文无法生成正确的Span
	opentracing.SetGlobalTracer(tracer)
	xxxmicro_opentracing.GlobalTracerWrapper().Wrap(tracer)
	logger.Info("create jaeger tracer success")
	return
}
