package model

import "sync/atomic"

type Metrics struct {
	operations atomic.Uint64
	incidents  atomic.Uint64
	trips      atomic.Uint64
	dropped    atomic.Uint64
}

func (m *Metrics) CountOperation() { m.operations.Add(1) }
func (m *Metrics) CountIncident()  { m.incidents.Add(1) }
func (m *Metrics) CountTrip()      { m.trips.Add(1) }
func (m *Metrics) CountDropped()   { m.dropped.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{
		"operations": m.operations.Load(),
		"incidents":  m.incidents.Load(),
		"trips":      m.trips.Load(),
		"dropped":    m.dropped.Load(),
	}
}
