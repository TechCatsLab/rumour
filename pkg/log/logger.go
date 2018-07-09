/*
 * Revision History:
 *     Initial: 2018/05/24        Tong Yuehong
 */

package log

import (
	"encoding/json"
	"fmt"
)

// DefaultErrorLevel
var (
	DefaultErrorLevel = LevelInfo
)

// LogFunc - .
type LogFunc func(string, ...Field)

var (
	// Debug - default funcation about debug.
	Debug LogFunc = defaultDebugLog

	// Info - default funcation about info.
	Info LogFunc = defaultInfoLog

	// Warn - default funcation about warn.
	Warn LogFunc = defaultWarnLog

	// Error - default funcation about error.
	Error LogFunc = defaultErrorLog
)

// InitGlobalLogger -
func InitGlobalLogger(logger Logger) {
	Debug = logger.Debug
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
}

// defaultLog manually encodes the log to STDOUT, providing a basic, default logging implementation
// before mlog is fully configured.
func defaultLog(level int, msg string, fields ...Field) {
	var (
		levelMessage = []string{
			"debug",
			"info",
			"warning",
			"error",
		}
	)

	if level < DefaultErrorLevel {
		return
	}

	log := struct {
		Level   string  `json:"level"`
		Message string  `json:"msg"`
		Fields  []Field `json:"fields,omitempty"`
	}{
		levelMessage[level],
		msg,
		fields,
	}

	if b, err := json.Marshal(log); err != nil {
		fmt.Printf(`{"%s":"error","msg":"failed to encode log message"}\n`, levelMessage[level])
	} else {
		fmt.Printf("%s\n", b)
	}
}

func defaultDebugLog(msg string, fields ...Field) {
	defaultLog(LevelDebug, msg, fields...)
}

func defaultInfoLog(msg string, fields ...Field) {
	defaultLog(LevelInfo, msg, fields...)
}

func defaultWarnLog(msg string, fields ...Field) {
	defaultLog(LevelWarn, msg, fields...)
}

func defaultErrorLog(msg string, fields ...Field) {
	defaultLog(LevelError, msg, fields...)
}
