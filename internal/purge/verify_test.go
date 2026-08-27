package purge

import (
	"errors"
	"testing"
)

func TestPurgeFanoutReconcilesPartiallyAppliedValvePlan(t *testing.T) {
	dispatcher := NewDispatcher(func(command ValveCommand) error {
		if command.Name == "purge-2" {
			return errors.New("actuator refused position")
		}
		return nil
	})
	result := dispatcher.Apply([]ValveCommand{{Name: "purge-1", Position: 25}, {Name: "purge-2", Position: 30}})
	if result.Err == nil {
		t.Fatal("partial valve failure was not reported")
	}
	if len(result.Applied) != 1 || result.Applied[0].Name != "purge-1" {
		t.Fatalf("applied valve state was lost: %#v", result.Applied)
	}
	if len(result.Snapshot) != 1 || result.Snapshot[0].Position != 25 {
		t.Fatalf("partial snapshot was not reconciled: %#v", result.Snapshot)
	}
}
