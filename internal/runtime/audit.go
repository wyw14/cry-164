package runtime

import (
	"fmt"
	"github.com/wyw14/cry-164/internal/analyzer"
	"github.com/wyw14/cry-164/internal/campaign"
	"github.com/wyw14/cry-164/internal/compressor"
	"github.com/wyw14/cry-164/internal/converter"
	"github.com/wyw14/cry-164/internal/cooling"
	"github.com/wyw14/cry-164/internal/feed"
	"github.com/wyw14/cry-164/internal/interlock"
	"github.com/wyw14/cry-164/internal/journal"
	"github.com/wyw14/cry-164/internal/model"
	"github.com/wyw14/cry-164/internal/purge"
	"github.com/wyw14/cry-164/internal/separator"
	"github.com/wyw14/cry-164/internal/steam"
)

func plantGraph() (*model.PlantGraph, error) {
	plans := []model.ComponentPlan{
		campaign.Plan(),
		feed.Plan(),
		compressor.Plan(),
		converter.Plan(),
		steam.Plan(),
		cooling.Plan(),
		separator.Plan(),
		purge.Plan(),
		analyzer.Plan(),
		interlock.Plan(),
		journal.Plan(),
	}
	graph, err := model.NewPlantGraph(plans...)
	if err != nil {
		return nil, fmt.Errorf("build plant graph: %w", err)
	}
	if len(graph.Components()) != len(plans) {
		return nil, fmt.Errorf("component graph lost plans")
	}
	if len(graph.Operations()) < 100 {
		return nil, fmt.Errorf("component graph is incomplete")
	}
	if len(graph.Equipment()) < 10 {
		return nil, fmt.Errorf("equipment graph is incomplete")
	}
	return graph, nil
}
