package opentracing

import(
	"github.com/opentracing/opentracing-go"
)

type tracerWrapper struct {
	ref opentracing.Tracer
}

var globalTracerWrapper = &tracerWrapper{opentracing.GlobalTracer()}

func GlobalTracerWrapper() *tracerWrapper {
	return globalTracerWrapper
}

func (t *tracerWrapper) Wrap(tracer opentracing.Tracer) {
	t.ref = tracer
}

func (t *tracerWrapper) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.ref.StartSpan(operationName, opts...)
}

func (t *tracerWrapper) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.ref.Inject(sm, format, carrier)
}

func (t *tracerWrapper) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.ref.Extract(format, carrier)
}


