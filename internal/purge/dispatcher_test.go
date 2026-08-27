package purge

import (
	"errors"
	"strings"
	"testing"
)

// appliedNames lists the valves recorded in Applied, in order.
func appliedNames(result Result) []string {
	names := make([]string, 0, len(result.Applied))
	for _, state := range result.Applied {
		names = append(names, state.Name)
	}
	return names
}

// snapshotPosition returns the recorded position for a valve, or -1 if absent.
func snapshotPosition(result Result, valve string) float64 {
	for _, state := range result.Snapshot {
		if state.Name == valve {
			return state.Position
		}
	}
	return -1
}

// TestApplyRejectsSecondValve reproduces the reported incident: when the second
// valve rejects its new position, the first valve has already moved physically.
// The committed snapshot must reflect that partial reality (valve 1 moved,
// valve 2 kept its prior position) so the page shows truth rather than a stale
// "old setting", and the next round computes from the actual valve positions.
func TestApplyRejectsSecondValve(t *testing.T) {
	rejected := errors.New("rejected")
	// failing holds the valve name the actuator rejects; empty means accept.
	failing := ""
	dispatcher := NewDispatcher(func(command ValveCommand) error {
		if failing != "" && command.Name == failing {
			return rejected
		}
		return nil
	})

	// Round 1 establishes the old committed setting.
	first := dispatcher.Apply([]ValveCommand{
		{Name: "purge-1", Position: 10},
		{Name: "purge-2", Position: 12},
	})
	if first.Err != nil {
		t.Fatalf("round 1 should succeed, got %v", first.Err)
	}

	// Round 2: purge-1 accepts and physically moves, purge-2 rejects.
	failing = "purge-2"
	failed := dispatcher.Apply([]ValveCommand{
		{Name: "purge-1", Position: 20},
		{Name: "purge-2", Position: 25},
	})
	if failed.Err == nil {
		t.Fatal("round 2 should fail with the rejected valve's error")
	}
	if !errors.Is(failed.Err, rejected) {
		t.Fatalf("error should wrap the actuator cause, got %v", failed.Err)
	}
	if got := appliedNames(failed); len(got) != 1 || got[0] != "purge-1" {
		t.Fatalf("Applied should record only the valve that moved, got %v", got)
	}
	if !strings.Contains(failed.Err.Error(), "purge-1=20") {
		t.Fatalf("error should name the valve that already moved, got %q", failed.Err.Error())
	}
	// Snapshot reflects the partial physical truth, not the old setting.
	if position := snapshotPosition(failed, "purge-1"); position != 20 {
		t.Fatalf("snapshot should reflect that purge-1 moved to 20, got %v", position)
	}
	if position := snapshotPosition(failed, "purge-2"); position != 12 {
		t.Fatalf("snapshot should keep purge-2 at its prior position 12, got %v", position)
	}

	// The next round computes from the actual positions: purge-1 starts at 20,
	// purge-2 at 12, not from the stale old setting of 10/12.
	failing = "" // accept all again
	next := dispatcher.Apply([]ValveCommand{
		{Name: "purge-1", Position: 30},
		{Name: "purge-2", Position: 35},
	})
	if next.Err != nil {
		t.Fatalf("round 3 should succeed, got %v", next.Err)
	}
	if position := snapshotPosition(next, "purge-1"); position != 30 {
		t.Fatalf("purge-1 should be 30, got %v", position)
	}
	if position := snapshotPosition(next, "purge-2"); position != 35 {
		t.Fatalf("purge-2 should be 35, got %v", position)
	}
	if len(next.Snapshot) != 2 {
		t.Fatalf("snapshot should hold one entry per valve, got %v", next.Snapshot)
	}
}

// TestApplyFirstValveRejected verifies that when nothing moved, the error does
// not claim any valve was applied and the snapshot keeps the prior position.
func TestApplyFirstValveRejected(t *testing.T) {
	stuck := false
	dispatcher := NewDispatcher(func(ValveCommand) error {
		if stuck {
			return errors.New("stuck")
		}
		return nil
	})
	if result := dispatcher.Apply([]ValveCommand{{Name: "purge-1", Position: 5}}); result.Err != nil {
		t.Fatalf("setup round should succeed, got %v", result.Err)
	}

	stuck = true
	result := dispatcher.Apply([]ValveCommand{{Name: "purge-1", Position: 9}})
	if result.Err == nil {
		t.Fatal("apply should fail")
	}
	if len(result.Applied) != 0 {
		t.Fatalf("no valve should be reported as applied, got %v", result.Applied)
	}
	if strings.Contains(result.Err.Error(), "already moved") {
		t.Fatalf("error should not mention moved valves, got %q", result.Err.Error())
	}
	if position := snapshotPosition(result, "purge-1"); position != 5 {
		t.Fatalf("snapshot should keep the prior position 5, got %v", position)
	}
}

// TestApplyUpdatesPerValve verifies positions are upserted per valve instead of
// accumulating duplicates across rounds.
func TestApplyUpdatesPerValve(t *testing.T) {
	dispatcher := NewDispatcher(func(ValveCommand) error { return nil })
	dispatcher.Apply([]ValveCommand{
		{Name: "purge-1", Position: 10},
		{Name: "purge-2", Position: 12},
	})
	second := dispatcher.Apply([]ValveCommand{
		{Name: "purge-1", Position: 20},
		{Name: "purge-2", Position: 22},
	})
	if len(second.Snapshot) != 2 {
		t.Fatalf("snapshot should have one entry per valve, got %v", second.Snapshot)
	}
	if position := snapshotPosition(second, "purge-1"); position != 20 {
		t.Fatalf("purge-1 should be 20, got %v", position)
	}
	if position := snapshotPosition(second, "purge-2"); position != 22 {
		t.Fatalf("purge-2 should be 22, got %v", position)
	}
}

// TestApplyEmptyPlan verifies an empty plan succeeds without changing state.
func TestApplyEmptyPlan(t *testing.T) {
	dispatcher := NewDispatcher(func(ValveCommand) error { return nil })
	result := dispatcher.Apply(nil)
	if result.Err != nil {
		t.Fatalf("empty plan should succeed, got %v", result.Err)
	}
	if len(result.Snapshot) != 0 {
		t.Fatalf("snapshot should be empty, got %v", result.Snapshot)
	}
}
