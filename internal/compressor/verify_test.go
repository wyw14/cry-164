package compressor

import (
	"testing"
	"time"
)

func TestLubeSoakConfigurationUsesSeconds(t *testing.T) {
	if got := LubricationDuration(45); got != 45*time.Second {
		t.Fatalf("45 seconds became %v", got)
	}
}
