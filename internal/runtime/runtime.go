package runtime

import (
	"context"
	"github.com/wyw14/cry-164/internal/analyzer"
	"github.com/wyw14/cry-164/internal/api"
	"github.com/wyw14/cry-164/internal/campaign"
	"github.com/wyw14/cry-164/internal/cooling"
	"github.com/wyw14/cry-164/internal/model"
	"github.com/wyw14/cry-164/internal/separator"
	"github.com/wyw14/cry-164/internal/service"
	"net/http"
)

type App struct {
	Server     *http.Server
	Runtime    *service.Runtime
	Operations *service.OperationService
	Plant      *service.Plant
}

func New(config model.Config) *App {
	state := campaign.NewState()
	tank := separator.NewTank()
	controller := cooling.NewController(tank)
	stream := &analyzer.Stream{}
	runtime := service.NewRuntime(state, controller, stream)
	plant := service.NewPlant(config, state, tank, controller, stream)
	graph, err := plantGraph()
	if err != nil {
		panic(err)
	}
	operations := service.NewOperationService(graph)
	handler := api.NewHandler(runtime, operations, plant)
	return &App{Server: &http.Server{Addr: ":" + itoa(config.Port), Handler: handler.Router()}, Runtime: runtime, Operations: operations, Plant: plant}
}
func (a *App) Start(ctx context.Context) error    { return a.Runtime.Start(ctx) }
func (a *App) Shutdown(ctx context.Context) error { return a.Server.Shutdown(ctx) }
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return digits
}
