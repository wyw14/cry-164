package analyzer

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestAnalyzerInitializationCanRecoverAfterTransientFailure(t *testing.T) {
	var calls atomic.Int32
	client := NewClient(func(context.Context) error {
		if calls.Add(1) == 1 {
			return errors.New("connection refused")
		}
		return nil
	})
	if err := client.Read(context.Background()); err == nil {
		t.Fatal("first handshake should report the transient failure")
	}
	if err := client.Read(context.Background()); err != nil {
		t.Fatalf("client did not recover after the device became available: %v", err)
	}
	if calls.Load() != 2 || !client.Ready() {
		t.Fatalf("unexpected recovered client state: calls=%d ready=%v", calls.Load(), client.Ready())
	}
}
