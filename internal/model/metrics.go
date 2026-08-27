package model

import "sync/atomic"

type Metrics struct {
	operations atomic.Uint64
	incidents  atomic.Uint64
}

func (m *Metrics) CountOperation() { m.operations.Add(1) }
func (m *Metrics) CountIncident()  { m.incidents.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"operations": m.operations.Load(), "incidents": m.incidents.Load()}
}
