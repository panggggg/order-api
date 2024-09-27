package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	logger *zap.Logger
}

type ZapConfig struct {
	Debug bool `env:"DEBUG" envDefault:"false"`
}

type CorrelationIdType string

const RequestIDKey CorrelationIdType = "requestID"

func NewLoggerZap(config *ZapConfig) (Logger, error) {

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.MessageKey = "message"
	loggerConfig.EncoderConfig.StacktraceKey = ""

	if config.Debug {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, _ := loggerConfig.Build(zap.AddCallerSkip(1))

	return &LoggerZap{
		logger: logger,
	}, nil
}

func fieldsToZap(ctx context.Context, fields []Field) []zap.Field {
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, Field{Key: "requestID", Value: ctx.Value(RequestIDKey)})
	}

	zapFields := make([]zap.Field, len(fields))

	for i, field := range fields {
		switch val := field.Value.(type) {
		case error:
			zapFields[i] = zap.Error(val)
		case string:
			zapFields[i] = zap.String(field.Key, val)
		case int:
			zapFields[i] = zap.Int(field.Key, val)
		case int64:
			zapFields[i] = zap.Int64(field.Key, val)
		case float64:
			zapFields[i] = zap.Float64(field.Key, val)
		case bool:
			zapFields[i] = zap.Bool(field.Key, val)
		case []byte:
			zapFields[i] = zap.ByteString(field.Key, val)
		case []string:
			zapFields[i] = zap.Strings(field.Key, val)
		default:
			zapFields[i] = zap.Any(field.Key, val)
		}
	}

	return zapFields
}

// implement the Logger interface

func (l *LoggerZap) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.logger.Fatal(msg, fieldsToZap(ctx, fields)...)
}

func (l *LoggerZap) Error(ctx context.Context, msg string, fields ...Field) {
	l.logger.Error(msg, fieldsToZap(ctx, fields)...)
}

func (l *LoggerZap) Debug(ctx context.Context, msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToZap(ctx, fields)...)
}

func (l *LoggerZap) Info(ctx context.Context, msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToZap(ctx, fields)...)
}

func (l *LoggerZap) Warn(ctx context.Context, msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToZap(ctx, fields)...)
}

func (l *LoggerZap) Panic(ctx context.Context, msg string, fields ...Field) {
	l.logger.Panic(msg, fieldsToZap(ctx, fields)...)
}
