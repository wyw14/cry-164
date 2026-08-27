package purge

import "context"

type EvidenceBarrier struct{}

func (EvidenceBarrier) Wait(ctx context.Context, sources []<-chan string) error {
	active := make([]<-chan string, 0, len(sources))
	for _, source := range sources {
		active = append(active, source)
	}
	for _, source := range active {
		select {
		case <-source:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
