package zerolog

import (
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	ZeroLogger.Error("hello")
}
