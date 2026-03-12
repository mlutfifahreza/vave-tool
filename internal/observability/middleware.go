package observability

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Middleware struct {
	logger          *zap.Logger
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
}

func NewMiddleware(logger *zap.Logger) (*Middleware, error) {
	meter := otel.Meter("vave-tool-api")

	requestCounter, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		logger:          logger,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	handler := otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx := r.Context()
			span := trace.SpanFromContext(ctx)
			traceID := span.SpanContext().TraceID().String()

			ctx = context.WithValue(ctx, "trace_id", traceID)
			r = r.WithContext(ctx)

			rr := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rr, r)

			duration := time.Since(start).Seconds()
			statusCode := rr.statusCode

			attrs := []attribute.KeyValue{
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.Int("status_code", statusCode),
			}

			m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
			m.requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))

			m.logger.Info("HTTP request",
				zap.String("trace_id", traceID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status_code", statusCode),
				zap.Float64("duration_seconds", duration),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		}),
		"HTTP Request",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)

	return handler
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	return rr.ResponseWriter.Write(b)
}
