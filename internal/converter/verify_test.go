package converter

import (
	"context"
	"testing"

	"github.com/wyw14/cry-164/internal/steam"
)

func TestConverterCleanupClosesEachOpenedBypassValve(t *testing.T) {
	first := &steam.Valve{}
	second := &steam.Valve{}
	err := RunHeating(context.Background(), []*steam.Valve{first, second}, 1)
	if err == nil {
		t.Fatal("heating failure was not reported")
	}
	if first.IsOpen() || second.IsOpen() {
		t.Fatalf("opened bypass valves were not closed: first=%v second=%v", first.IsOpen(), second.IsOpen())
	}
}
