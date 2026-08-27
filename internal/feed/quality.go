package feed

import "github.com/wyw14/cry-164/internal/model"

type QualityGate struct{ limit float64 }

func NewQualityGate(limit float64) QualityGate { return QualityGate{limit: limit} }
func (g QualityGate) Accept(reading model.Reading) bool {
	return reading.Hydrogen >= g.limit && reading.Nitrogen > 0
}
func (g QualityGate) Describe() string { return "feed-quality" }
