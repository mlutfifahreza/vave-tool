package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	dbCallDuration  metric.Float64Histogram
	dbCallCounter   metric.Int64Counter
	cacheHitCounter metric.Int64Counter
	cacheCounter    metric.Int64Counter
	redisOpDuration metric.Float64Histogram
	redisOpCounter  metric.Int64Counter
}

func InitMetrics() (*Metrics, error) {
	meter := otel.Meter("vave-tool-db")

	dbCallDuration, err := meter.Float64Histogram(
		"db_call_duration_seconds",
		metric.WithDescription("Database operation duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	dbCallCounter, err := meter.Int64Counter(
		"db_calls_total",
		metric.WithDescription("Total number of database calls"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		return nil, err
	}

	cacheHitCounter, err := meter.Int64Counter(
		"cache_hits_total",
		metric.WithDescription("Total number of cache hits and misses"),
		metric.WithUnit("{hit}"),
	)
	if err != nil {
		return nil, err
	}

	cacheCounter, err := meter.Int64Counter(
		"cache_operations_total",
		metric.WithDescription("Total number of cache operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	redisOpDuration, err := meter.Float64Histogram(
		"redis_operation_duration_seconds",
		metric.WithDescription("Redis operation duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	redisOpCounter, err := meter.Int64Counter(
		"redis_operations_total",
		metric.WithDescription("Total number of Redis operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		dbCallDuration:  dbCallDuration,
		dbCallCounter:   dbCallCounter,
		cacheHitCounter: cacheHitCounter,
		cacheCounter:    cacheCounter,
		redisOpDuration: redisOpDuration,
		redisOpCounter:  redisOpCounter,
	}, nil
}

func (m *Metrics) RecordDBCall(ctx context.Context, operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("status", status),
	}

	m.dbCallCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.dbCallDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *Metrics) RecordCacheAccess(ctx context.Context, operation string, hit bool) {
	hitStatus := "miss"
	if hit {
		hitStatus = "hit"
	}

	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("status", hitStatus),
	}

	m.cacheHitCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *Metrics) RecordRedisOp(ctx context.Context, operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("status", status),
	}

	m.redisOpCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.redisOpDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

var globalMetrics *Metrics

func SetGlobalMetrics(m *Metrics) {
	globalMetrics = m
}

func GetMetrics() *Metrics {
	return globalMetrics
}
