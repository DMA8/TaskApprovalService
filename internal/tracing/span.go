package tracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan("ping-receive", ext.RPCServerOption(spanCtx))
}