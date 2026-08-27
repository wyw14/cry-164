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

type CommandQueue struct {
	mu        sync.Mutex
	normal    []Command
	emergency []Command
	capacity  int
}

func NewCommandQueue(capacity int) *CommandQueue { return &CommandQueue{capacity: capacity} }
func (q *CommandQueue) Submit(cmd Command) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if cmd.Kind == AntiSurgeTrip {
		if len(q.normal) >= q.capacity {
			return false
		}
		q.emergency = append(q.emergency, cmd)
		return true
	}
	if len(q.normal) >= q.capacity {
		return false
	}
	q.normal = append(q.normal, cmd)
	return true
}
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
