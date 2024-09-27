package log

import "context"

// Field is a key-value pair that is used to add structured context to a log message.
type Field struct {
	Key   string
	Value interface{}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Bytes(key string, value []byte) Field {
	return Field{Key: key, Value: value}
}

func Strings(key string, value []string) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

// Fields is a slice of Field.
type Fields []Field

// With returns a new Fields with the given key-value pair appended.
func (f Fields) With(key string, value interface{}) Fields {
	return append(f, Field{Key: key, Value: value})
}

// Logger is the interface that wraps the basic logging methods.
type Logger interface {

	// Panic logs a message at panic level.
	Panic(ctx context.Context, msg string, fields ...Field)

	// Fatal logs a message at fatal level.
	Fatal(ctx context.Context, msg string, fields ...Field)

	// Error logs a message at error level.
	Error(ctx context.Context, msg string, fields ...Field)

	// Warn logs a message at warn level.
	Warn(ctx context.Context, msg string, fields ...Field)

	// Info logs a message at info level.
	Info(ctx context.Context, msg string, fields ...Field)

	// Debug logs a message at debug level.
	Debug(ctx context.Context, msg string, fields ...Field)
}
