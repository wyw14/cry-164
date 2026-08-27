package compressor

import (
	"encoding/json"
	"errors"
	"time"
)

type Frame struct {
	Pressure    float64 `json:"pressure"`
	Temperature float64 `json:"temperature"`
	Trip        bool    `json:"trip"`
}
type HotDischargeError struct{ Temperature float64 }

func (e *HotDischargeError) Error() string { return "hot discharge" }

type OtherError struct{ Message string }

func (e *OtherError) Error() string { return e.Message }

func DecodeFrame(payload []byte, frame *Frame) error {
	*frame = Frame{}
	return json.Unmarshal(payload, frame)
}

func LubricationDuration(seconds int) time.Duration { return time.Duration(seconds) * time.Second }

func Classify(err error) string {
	var hot *HotDischargeError
	if errors.As(err, &hot) {
		return "emergency-unload"
	}
	return "degraded"
}
