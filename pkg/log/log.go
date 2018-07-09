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

// Field - add a key-value pair to a logger's context.
type Field = zapcore.Field

var (
	// Int64 return a field with the given key and value.
	Int64 = zap.Int64

	// Int return a field with the given key and value.
	Int = zap.Int

	// String return a field with the given key and value.
	String = zap.String

	// Any return a field with a key and an arbitrary value.
	Any = zap.Any

	// Err return a field by using NamedError.
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
	// LevelDebug - represent the level of Debug.
	LevelDebug = iota
	// LevelInfo - represent the level of Info.
	LevelInfo
	// LevelWarn - represent the level of Warn.
	LevelWarn
	// LevelError - represent the level of Error.
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
