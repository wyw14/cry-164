package compressor

import "testing"

func TestAntiSurgeTripCannotBeDroppedByBusyCommandQueue(t *testing.T) {
	queue := NewCommandQueue(1)
	if !queue.Submit(Command{Kind: LoadAdjust, Value: 95}) {
		t.Fatal("normal command should fill the queue")
	}
	if !queue.Submit(Command{Kind: AntiSurgeTrip}) {
		t.Fatal("anti-surge trip was rejected by normal queue pressure")
	}
	command, ok := queue.Next()
	if !ok || command.Kind != AntiSurgeTrip {
		t.Fatalf("expected emergency command first, got %#v, %v", command, ok)
	}
}
