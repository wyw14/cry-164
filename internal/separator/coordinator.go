package separator

import "context"

type Coordinator struct {
	tank *Tank
}

func NewCoordinator(tank *Tank) *Coordinator {
	return &Coordinator{tank: tank}
}

func (c *Coordinator) AuthorizeDrain(ctx context.Context) error {
	return c.tank.WaitUntilReady(ctx)
}
