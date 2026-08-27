package analyzer

import (
	"context"
	"errors"
	"testing"
)

func TestAnalyzerRetryCancelsEachAttemptPromptly(t *testing.T) {
	poller := NewPoller(func(context.Context) error { return errors.New("offline") })
	if err := poller.Poll(context.Background(), 12); err == nil {
		t.Fatal("offline analyzer should exhaust retries")
	}
	if got := poller.MaxOutstanding(); got != 1 {
		t.Fatalf("attempt resources accumulated across retries: max outstanding=%d", got)
	}
}
