package purge

import (
	"context"
	"testing"
	"time"
)

func TestDisabledPurgeSourceDoesNotBlockEvidenceBarrier(t *testing.T) {
	ready := make(chan string, 1)
	ready <- "current-cycle"
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	if err := (EvidenceBarrier{}).Wait(ctx, []<-chan string{ready, nil}); err != nil {
		t.Fatalf("disabled source blocked the evidence barrier: %v", err)
	}
}
