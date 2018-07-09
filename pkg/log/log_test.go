package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	Debug("debug")
	Info("info")

	InitGlobalLogger(*NewLogger(1))

	Info("info")
	Debug("debug")
}
