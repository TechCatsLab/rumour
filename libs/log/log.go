/*
 * Revision History:
 *     Initial: 2018/05/24        Tong Yuehong
 */

package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field - Type and function aliases from zap to limit the libraries scope into MM code
type Field = zapcore.Field

var (
	// Int64 constructs a field with the given key and value.
	Int64 = zap.Int64

	// Int constructs a field with the given key and value.
	Int = zap.Int

	// String constructs a field with the given key and value.
	String = zap.String

	// Any takes a key and an arbitrary value and chooses the best way to represent
	// them as a field, falling back to a reflection-based approach only if
	// necessary.
	Any = zap.Any

	// Err is shorthand for the common idiom NamedError("error", err).
	Err = zap.Error

	levelTranformer = []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}
)

// Logger -
type Logger struct {
	zap   *zap.Logger
	level zapcore.Level
}

const (
	// LevelDebug -
	LevelDebug = iota
	// LevelInfo -
	LevelInfo
	// LevelWarn -
	LevelWarn
	// LevelError -
	LevelError
)

func getZapLevel(level int) zapcore.Level {
	return levelTranformer[level]
}

func makeEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// NewLogger - new a Logger.
func NewLogger(level int) *Logger {
	l := getZapLevel(level)
	logger := &Logger{
		level: l,
	}

	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(makeEncoder(), writer, l)
	logger.zap = zap.New(core,
		zap.AddCallerSkip(2),
		zap.AddCaller())

	return logger
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(message string, fields ...Field) {
	l.zap.Debug(message, fields...)
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(message string, fields ...Field) {
	l.zap.Info(message, fields...)
}

// Warn logs a message at WarnLevel
func (l *Logger) Warn(message string, fields ...Field) {
	l.zap.Warn(message, fields...)
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(message string, fields ...Field) {
	l.zap.Error(message, fields...)
}
