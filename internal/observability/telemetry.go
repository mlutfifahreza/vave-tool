package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Telemetry struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *metric.MeterProvider
	MetricsHandler http.Handler
	Logger         *zap.Logger
}

func InitTelemetry(serviceName, serviceVersion, otelEndpoint string) (*Telemetry, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider, err := initTracer(ctx, res, otelEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	meterProvider, err := initMetrics(res)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	logger := initLogger(serviceName)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider.Provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider.Provider,
		MetricsHandler: meterProvider.Handler,
		Logger:         logger,
	}, nil
}

func initTracer(ctx context.Context, res *resource.Resource, endpoint string) (*sdktrace.TracerProvider, error) {
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tracerProvider, nil
}

func initMetrics(res *resource.Resource) (*struct {
	Provider *metric.MeterProvider
	Handler  http.Handler
}, error) {
	registry := prometheus.NewRegistry()

	exporter, err := promexporter.New(
		promexporter.WithRegisterer(registry),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	histogramView := metric.NewView(
		metric.Instrument{Name: "http_request_duration_seconds"},
		metric.Stream{
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{
					0.001, 0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0, 7.5, 10.0,
				},
			},
		},
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exporter),
		metric.WithView(histogramView),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return &struct {
		Provider *metric.MeterProvider
		Handler  http.Handler
	}{
		Provider: meterProvider,
		Handler:  handler,
	}, nil
}

func initLogger(serviceName string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.InitialFields = map[string]interface{}{
		"service_name": serviceName,
	}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	return logger
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := t.TracerProvider.Shutdown(ctx); err != nil {
		t.Logger.Error("Error shutting down tracer provider", zap.Error(err))
		return err
	}

	if err := t.MeterProvider.Shutdown(ctx); err != nil {
		t.Logger.Error("Error shutting down meter provider", zap.Error(err))
		return err
	}

	if err := t.Logger.Sync(); err != nil {
		return err
	}

	return nil
}
