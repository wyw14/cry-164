package service

import "context"

type Task func(context.Context) error

func StopComponents(ctx context.Context, analyze Task, cool Task) error {
	analysisCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	analysisErr := make(chan error, 1)
	go func() { analysisErr <- analyze(analysisCtx) }()
	coolingErr := cool(ctx)
	if coolingErr != nil {
		return coolingErr
	}
	return <-analysisErr
}
