package compressor

import "sync"

type CommandKind string

const (
	LoadAdjust    CommandKind = "load-adjust"
	AntiSurgeTrip CommandKind = "anti-surge-trip"
)

type Command struct {
	Kind  CommandKind
	Value float64
}

// CommandQueue separates emergency (anti-surge trip) commands from normal
// load-adjust commands. Emergency commands never get dropped behind a full
// normal lane — a protection trip must always be admitted, so it owns a
// dedicated emergency lane and is drained with priority.
type CommandQueue struct {
	mu        sync.Mutex
	normal    []Command
	emergency []Command
	capacity  int
	dropped   int
}

func NewCommandQueue(capacity int) *CommandQueue { return &CommandQueue{capacity: capacity} }

// Submit enqueues a command. Emergency (anti-surge trip) commands are admitted
// to the emergency lane regardless of normal-lane occupancy; they are only
// dropped if the emergency lane itself is saturated, which is recorded in
// Dropped(). Normal commands are dropped (and counted) when the normal lane is
// full. Returns false when the command was dropped.
func (q *CommandQueue) Submit(cmd Command) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if cmd.Kind == AntiSurgeTrip {
		if len(q.emergency) >= q.capacity {
			q.dropped++
			return false
		}
		q.emergency = append(q.emergency, cmd)
		return true
	}
	if len(q.normal) >= q.capacity {
		q.dropped++
		return false
	}
	q.normal = append(q.normal, cmd)
	return true
}

// Next drains the emergency lane with priority before the normal lane.
func (q *CommandQueue) Next() (Command, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.emergency) > 0 {
		cmd := q.emergency[0]
		q.emergency = q.emergency[1:]
		return cmd, true
	}
	if len(q.normal) > 0 {
		cmd := q.normal[0]
		q.normal = q.normal[1:]
		return cmd, true
	}
	return Command{}, false
}

// Dropped reports the number of commands that have been rejected because their
// lane was full, including emergency trips that could not be admitted.
func (q *CommandQueue) Dropped() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.dropped
}
