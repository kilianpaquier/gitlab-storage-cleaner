package engine

import (
	"context"
	"fmt"
	"io"

	"github.com/panjf2000/ants/v2"
)

// Logger represents a basic interface with minimal methods used during artifacts cleanup feature.
type Logger interface {
	ants.Logger

	// Info logs a message at level INFO.
	Info(msg any, keyvals ...any)

	// Error logs a message at level ERROR.
	Error(msg any, keyvals ...any)

	// Warn logs a message at level WARN.
	Warn(msg any, keyvals ...any)

	// Debug logs a message at level DEBUG.
	Debug(msg any, keyvals ...any)
}

type loggerKeyType string

// LoggerKey is the context key for the logger.
const LoggerKey loggerKeyType = "logger"

// GetLogger returns the context logger.
//
// By default it will a noop logger, but it can be set with WithLogger run option.
func GetLogger(ctx context.Context) Logger {
	log, ok := ctx.Value(LoggerKey).(Logger)
	if !ok {
		return &noopLogger{}
	}
	return log
}

type noopLogger struct{}

var _ Logger = &noopLogger{} // ensure interface is implemented

// Printf does nothing.
func (*noopLogger) Printf(string, ...any) {}

// Debug does nothing.
func (*noopLogger) Debug(any, ...any) {}

// Error does nothing.
func (*noopLogger) Error(any, ...any) {}

// Info does nothing.
func (*noopLogger) Info(any, ...any) {}

// Warn does nothing.
func (*noopLogger) Warn(any, ...any) {}

type testLogger struct{ writer io.Writer }

var _ Logger = (*testLogger)(nil) // ensure interface is implemented

// NewTestLogger creates a new logger with the input writer.
//
// This logger is expected to be used in tests.
// In no way it should be used in production since it's unoptimized.
func NewTestLogger(writer io.Writer) Logger {
	return &testLogger{writer: writer}
}

// Debug implements shared.Logger.
func (b *testLogger) Debug(msg any, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Error implements shared.Logger.
func (b *testLogger) Error(msg any, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Info implements shared.Logger.
func (b *testLogger) Info(msg any, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Printf implements shared.Logger.
func (b *testLogger) Printf(format string, args ...any) {
	b.writer.Write([]byte(fmt.Sprintf(format, args...)))
}

// Warn implements shared.Logger.
func (b *testLogger) Warn(msg any, keyvals ...any) {
	b.print(msg, keyvals...)
}

func (b *testLogger) print(msg any, keyvals ...any) {
	b.writer.Write([]byte(fmt.Sprint(msg)))
	for i := 0; i < len(keyvals); i += 2 {
		b.writer.Write([]byte(fmt.Sprintf(" %s=%v", keyvals[i], keyvals[i+1])))
	}
	b.writer.Write([]byte("\n"))
}
