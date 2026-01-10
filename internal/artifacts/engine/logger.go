package engine

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/panjf2000/ants/v2"
)

// Logger represents a basic interface with minimal methods used during artifacts cleanup feature.
type Logger interface {
	ants.Logger

	// Info logs a message at level INFO.
	Info(msg string, keyvals ...any)

	// Error logs a message at level ERROR.
	Error(msg string, keyvals ...any)

	// Warn logs a message at level WARN.
	Warn(msg string, keyvals ...any)

	// Debug logs a message at level DEBUG.
	Debug(msg string, keyvals ...any)
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
func (*noopLogger) Debug(string, ...any) {}

// Error does nothing.
func (*noopLogger) Error(string, ...any) {}

// Info does nothing.
func (*noopLogger) Info(string, ...any) {}

// Warn does nothing.
func (*noopLogger) Warn(string, ...any) {}

type testLogger struct{ writer io.Writer }

var _ Logger = (*testLogger)(nil) // ensure interface is implemented

// NewTestLogger creates a new logger with the input writer.
//
// This logger is expected to be used in tests.
// In no way it should be used in production since it's unoptimized.
func NewTestLogger(writer io.Writer) Logger {
	return &testLogger{writer: writer}
}

// Debug implements Logger.
func (b *testLogger) Debug(msg string, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Error implements Logger.
func (b *testLogger) Error(msg string, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Info implements Logger.
func (b *testLogger) Info(msg string, keyvals ...any) {
	b.print(msg, keyvals...)
}

// Printf implements Logger.
func (b *testLogger) Printf(format string, args ...any) {
	b.writer.Write(fmt.Appendf(nil, format, args...))
}

// Warn implements Logger.
func (b *testLogger) Warn(msg string, keyvals ...any) {
	b.print(msg, keyvals...)
}

func (b *testLogger) print(msg any, keyvals ...any) {
	b.writer.Write(fmt.Append(nil, msg))
	for i := 0; i < len(keyvals); i += 2 {
		b.writer.Write(fmt.Appendf(nil, " %s=%v", keyvals[i], keyvals[i+1]))
	}
	b.writer.Write([]byte("\n"))
}

func NewSlogLogger(log *slog.Logger) Logger {
	return &slogLogger{log}
}

type slogLogger struct {
	log *slog.Logger
}

var _ Logger = (*slogLogger)(nil)

// Printf implements Logger.
func (s *slogLogger) Printf(format string, args ...any) {
	s.log.Info(fmt.Sprintf(format, args...))
}

// Debug implements Logger.
func (s *slogLogger) Debug(msg string, keyvals ...any) {
	s.log.Debug(msg, keyvals...)
}

// Error implements Logger.
func (s *slogLogger) Error(msg string, keyvals ...any) {
	s.log.Error(msg, keyvals...)
}

// Info implements Logger.
func (s *slogLogger) Info(msg string, keyvals ...any) {
	s.log.Info(msg, keyvals...)
}

// Warn implements Logger.
func (s *slogLogger) Warn(msg string, keyvals ...any) {
	s.log.Warn(msg, keyvals...)
}
