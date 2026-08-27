package converter

import "context"

type ShutdownTask func(context.Context) error

func RunShutdown(ctx context.Context, optional ShutdownTask, quench ShutdownTask) error {
	optionalCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	optionalResult := make(chan error, 1)
	go func() {
		optionalResult <- optional(optionalCtx)
	}()
	quenchErr := quench(ctx)
	<-optionalResult
	if quenchErr != nil {
		return quenchErr
	}
	return nil
}
