package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger(zapLogger *zap.Logger) *Logger {
	return &Logger{zap: zapLogger}
}

func (l *Logger) WithTraceID(ctx context.Context) *zap.Logger {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()

	if traceID != "" && traceID != "00000000000000000000000000000000" {
		return l.zap.With(zap.String("trace_id", traceID))
	}

	return l.zap
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithTraceID(ctx).Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithTraceID(ctx).Error(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithTraceID(ctx).Warn(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithTraceID(ctx).Debug(msg, fields...)
}

func StartSpan(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer := otel.Tracer("vave-tool-api")
	return tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

func RecordError(span trace.Span, err error, msg string) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, msg)
	}
}

func AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}
