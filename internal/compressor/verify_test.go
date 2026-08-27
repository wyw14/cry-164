package compressor

import "testing"

func TestSparseCompressorFrameDoesNotCarryPriorTripFlag(t *testing.T) {
	frame := Frame{Pressure: 100, Temperature: 410, Trip: true}
	if err := DecodeFrame([]byte(`{"pressure":118}`), &frame); err != nil {
		t.Fatalf("decode sparse frame: %v", err)
	}
	if frame.Trip {
		t.Fatal("sparse frame inherited the prior trip flag")
	}
	if frame.Pressure != 118 {
		t.Fatalf("new pressure was not decoded: %v", frame.Pressure)
	}
}
