package compressor

import (
	"fmt"
	"testing"
)

func TestWrappedHotDischargeErrorRemainsClassifiable(t *testing.T) {
	err := fmt.Errorf("compressor protection: %w", &HotDischargeError{Temperature: 485})
	if got := Classify(err); got != "emergency-unload" {
		t.Fatalf("wrapped hot-discharge error was classified as %q", got)
	}
}
