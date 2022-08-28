package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	namespace = "team31"
	subsystem = "tasks"
)

var tracer trace.Tracer

func TracingMiddleware(next http.Handler) http.Handler {
	newTracerProvider()
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "tracingMiddleware")
		next.ServeHTTP(w, r.WithContext(ctx))
		defer span.End()
		log.Println("tracing middleware ok")
	}))
}

func newTracerProvider() (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://jaeger-instance-collector.observability:14268/api/traces")),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("team31_tasks"),
		)))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer = tp.Tracer("team31_tasks")
	log.Println("tracer is set")
	return tp, nil
}


func Logger(l *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				Entry := logging.Entry{
					Service:      "task",
					Method:       r.Method,
					Url:          r.URL.Path,
					Query:        r.URL.RawQuery,
					RemoteIP:     r.RemoteAddr,
					Status:       ww.Status(),
					Size:         ww.BytesWritten(),
					ReceivedTime: start,
					Duration:     time.Since(start),
					ServerIP:     r.Host,
					UserAgent:    r.Header.Get("User-Agent"),
					RequestId:    GetReqID(r.Context()),
				}
				l.Info().Msgf("%+v", Entry)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Prometheus() func(http.Handler) http.Handler {
	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"path", "code"},
	)

	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "latency",
			Help:      "Request-response latency histogram",
			Buckets: []float64{0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1,
				1.2, 1.4, 1.6, 1.8, 2, 2.3, 2.6, 2.9, 3.2, 3.6, 3.9,
				4.4, 4.9, 5.4, 5.9, 7, 8, 9, 10, 11},
		},
		[]string{"path"},
	)

	prometheus.Register(requestsCounter)
	prometheus.Register(latencyHistogram)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9000", nil)

	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				status := strconv.Itoa(ww.Status())
				requestsCounter.WithLabelValues(r.URL.Path, status).Inc()
				latencyHistogram.WithLabelValues(r.URL.Path).
					Observe(time.Since(start).Seconds())
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

