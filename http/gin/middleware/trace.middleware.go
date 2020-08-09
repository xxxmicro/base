package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type TraceMiddleware struct {
}

func (m *TraceMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var parentSpan opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			parentSpan = opentracing.StartSpan(ctx.Request.URL.Path)
		} else {
			parentSpan = opentracing.StartSpan(
				ctx.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer parentSpan.Finish()
		
		ctx.Set("ParentSpanContext", parentSpan.Context())
		ctx.Next()
	}
}

func NewTraceMiddleware() *TraceMiddleware {
	return &TraceMiddleware{}
}