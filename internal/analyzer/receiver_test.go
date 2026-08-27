package analyzer

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestClientReadCachesSuccess(t *testing.T) {
	var calls int32
	handshake := func(context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}
	c := NewClient(handshake)

	for i := 0; i < 5; i++ {
		if err := c.Read(context.Background()); err != nil {
			t.Fatalf("read %d: %v", i, err)
		}
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("handshake should run once on success, got %d", got)
	}
	if !c.Ready() {
		t.Fatal("client should report ready after successful handshake")
	}
}

// TestClientReadRecoversAfterTransientFailure reproduces the regression: a
// device that was offline at startup (analyzer started a minute late) returns
// connection refused on the first handshake. The client must not pin that
// failure forever; once the device is back online the next read must retry the
// handshake and succeed.
func TestClientReadRecoversAfterTransientFailure(t *testing.T) {
	var calls int32
	handshake := func(context.Context) error {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			return errors.New("connection refused")
		}
		return nil
	}
	c := NewClient(handshake)

	// First read: device still booting, handshake fails.
	err := c.Read(context.Background())
	if err == nil || err.Error() != "connection refused" {
		t.Fatalf("first read should fail with connection refused, got %v", err)
	}
	if c.Ready() {
		t.Fatal("client must not be ready after a failed handshake")
	}

	// Device comes back online: subsequent read must retry the handshake.
	if err := c.Read(context.Background()); err != nil {
		t.Fatalf("read after recovery should succeed, got %v", err)
	}
	if !c.Ready() {
		t.Fatal("client should report ready after recovered handshake")
	}

	// Further reads reuse the established session without re-handshaking.
	if err := c.Read(context.Background()); err != nil {
		t.Fatalf("steady-state read should succeed, got %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("handshake should run exactly twice (fail then success), got %d", got)
	}
}

// TestPollerRetriesHandshakeAfterFailure ensures the analyzer-retry path
// (Poller.Poll over client.Read) actually issues new connection attempts once
// the device is back online, instead of returning a cached error.
//
// Models the real outage: the analyzer is offline for a whole minute, so every
// handshake during that window fails. The first Poll (covering the offline
// window) exhausts its attempts and returns the error. Once the device is back
// online, the next Poll must issue a fresh handshake that succeeds.
func TestPollerRetriesHandshakeAfterFailure(t *testing.T) {
	var calls int32
	offline := atomic.Bool{}
	offline.Store(true)
	handshake := func(context.Context) error {
		atomic.AddInt32(&calls, 1)
		if offline.Load() {
			return errors.New("connection refused")
		}
		return nil
	}
	c := NewClient(handshake)
	poller := NewPoller(c.Read)

	// Device offline for the whole first poll: all attempts fail. The poller
	// must actually try the handshake on each attempt (no cached failure).
	if err := poller.Poll(context.Background(), 3); err == nil {
		t.Fatal("first poll should fail while device is offline")
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("offline poll should attempt 3 fresh handshakes, got %d", got)
	}

	// Device comes back online: the next poll must retry the handshake and
	// succeed rather than returning the cached failure.
	offline.Store(false)
	if err := poller.Poll(context.Background(), 3); err != nil {
		t.Fatalf("poll after recovery should succeed, got %v", err)
	}
	// First attempt of the recovered poll succeeds immediately.
	if got := atomic.LoadInt32(&calls); got != 4 {
		t.Fatalf("recovered poll should issue one fresh handshake, got %d total", got)
	}
}
