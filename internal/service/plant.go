package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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

type Plant struct {
	mu              sync.Mutex
	config          model.Config
	cycle           *campaign.State
	startup         *campaign.Coordinator
	permit          *feed.Permit
	machine         *compressor.Startup
	commands        *compressor.CommandQueue
	session         *compressor.Session
	analyzerClient  *analyzer.Client
	analyzerPoller  *analyzer.Poller
	analyzerState   *analyzer.State
	stream          *analyzer.Stream
	quality         feed.QualityGate
	steam           *steam.Recovery
	converter       *converter.Controller
	beds            *converter.BedState
	valves          []*steam.Valve
	tank            *separator.Tank
	separator       *separator.Coordinator
	cooling         *cooling.Controller
	telemetry       *cooling.Telemetry
	purgeDispatcher *purge.Dispatcher
	purgePlans      *purge.PlanStore
	lastPurgeAudit  purge.AuditEntry
	interlock       *interlock.State
	journal         *journal.Store
	journalPath     string
	metrics         model.Metrics
}

func NewPlant(config model.Config, cycle *campaign.State, tank *separator.Tank, coolingController *cooling.Controller, stream *analyzer.Stream) *Plant {
	permit := feed.NewPermit(120)
	machine := compressor.NewStartup(5 * time.Millisecond)
	client := analyzer.NewClient(func(context.Context) error { return nil })
	poller := analyzer.NewPoller(client.Read)
	recovery := steam.NewRecovery()
	journalPath := strings.TrimSpace(os.Getenv("AMMONIALOOP_JOURNAL"))
	if journalPath == "" {
		journalPath = filepath.Join(os.TempDir(), "cry-164-ammonialoop-incidents.json")
	}
	return &Plant{
		config:          config,
		cycle:           cycle,
		startup:         campaign.NewCoordinator(permit, machine, config),
		permit:          permit,
		machine:         machine,
		commands:        compressor.NewCommandQueue(16),
		session:         compressor.NewSession("recycle-compressor"),
		analyzerClient:  client,
		analyzerPoller:  poller,
		analyzerState:   &analyzer.State{},
		stream:          stream,
		quality:         feed.NewQualityGate(2.5),
		steam:           recovery,
		converter:       converter.NewController(stream, recovery),
		beds:            &converter.BedState{},
		valves:          []*steam.Valve{{}, {}, {}},
		tank:            tank,
		separator:       separator.NewCoordinator(tank),
		cooling:         coolingController,
		telemetry:       &cooling.Telemetry{},
		purgeDispatcher: purge.NewDispatcher(func(purge.ValveCommand) error { return nil }),
		purgePlans:      &purge.PlanStore{},
		interlock:       &interlock.State{},
		journal:         journal.NewStore(journalPath),
		journalPath:     journalPath,
	}
}

func (p *Plant) Execute(ctx context.Context, action, detail string) (any, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	action = strings.TrimSpace(action)
	detail = strings.TrimSpace(detail)
	switch action {
	case "startup":
		lubrication := p.config.LubeSoak()
		if detail != "" {
			value, err := time.ParseDuration(detail)
			if err != nil {
				return nil, fmt.Errorf("parse lubrication duration: %w", err)
			}
			lubrication = value
		}
		err := p.startup.Start(ctx, campaign.StartupPlan{Lubrication: lubrication, StageTimeout: p.config.StageTimeout})
		return map[string]any{"running": p.machine.Running()}, err
	case "stop-compressor":
		p.machine.Stop()
		return map[string]any{"running": p.machine.Running()}, nil
	case "load-adjust", "anti-surge-trip":
		kind := compressor.LoadAdjust
		if action == "anti-surge-trip" {
			kind = compressor.AntiSurgeTrip
		}
		value := 0.0
		if detail != "" {
			parsed, err := strconv.ParseFloat(detail, 64)
			if err != nil {
				return nil, fmt.Errorf("parse command value: %w", err)
			}
			value = parsed
		}
		if !p.commands.Submit(compressor.Command{Kind: kind, Value: value}) {
			p.metrics.CountDropped()
			return nil, errors.New("compressor command queue is full")
		}
		// An anti-surge trip is a latching protection: it must act on the
		// machine immediately, not merely sit in the command queue. Latch the
		// compressor trip, declare the interlock emergency, journal the trip so
		// the unit record carries the action, and force the production cycle
		// back to preparing so the campaign does not keep running while the
		// machine is tripped.
		if kind == compressor.AntiSurgeTrip {
			p.session.Trip()
			tripReason := fmt.Sprintf("anti-surge trip at load %.2f", value)
			p.interlock.SetEmergency(tripReason)
			incident := journal.IncidentFor(tripReason, "operator")
			if err := p.journal.Append(incident); err != nil {
				return nil, fmt.Errorf("journal anti-surge trip: %w", err)
			}
			p.metrics.CountIncident()
			p.metrics.CountTrip()
			// A surge trip can fire from any live cycle state, and the guarded
			// Advance() chain only moves a campaign forward, never back to
			// preparing. Reset the campaign to a fresh preparing cycle so
			// production does not keep running on a tripped machine.
			p.cycle.Reset()
			_, ok := p.commands.Next()
			return map[string]any{
				"accepted":  true,
				"dispatched": ok,
				"kind":      kind,
				"value":     value,
				"tripped":   true,
				"incident":  incident,
				"dropped":   p.commands.Dropped(),
			}, nil
		}
		command, ok := p.commands.Next()
		return map[string]any{"accepted": true, "dispatched": ok, "kind": command.Kind, "value": command.Value, "dropped": p.commands.Dropped()}, nil
	case "analyzer-sample":
		if err := p.analyzerClient.Read(ctx); err != nil {
			return nil, err
		}
		reading := model.Reading{Source: "loop-analyzer", Hydrogen: 3, Nitrogen: 1, Ammonia: 0.12, Sequence: p.analyzerState.Snapshot().Sequence + 1}
		p.analyzerState.Replace(reading)
		p.stream.Publish(reading)
		p.beds.Replace([]float64{reading.Ammonia, reading.Hydrogen, reading.Nitrogen})
		return map[string]any{"reading": reading, "accepted": p.quality.Accept(reading)}, nil
	case "analyzer-retry":
		err := p.analyzerPoller.Poll(ctx, 3)
		return map[string]any{"ready": p.analyzerClient.Ready(), "max_outstanding": p.analyzerPoller.MaxOutstanding()}, err
	case "decode-compressor-frame":
		var frame compressor.Frame
		if err := compressor.DecodeFrame([]byte(detail), &frame); err != nil {
			return nil, fmt.Errorf("decode compressor frame: %w", err)
		}
		p.session.Update(frame)
		return compressor.EquipmentState(p.session), nil
	case "lube-duration":
		seconds, err := strconv.Atoi(detail)
		if err != nil {
			return nil, fmt.Errorf("parse lube seconds: %w", err)
		}
		return map[string]any{"duration": compressor.LubricationDuration(seconds).String()}, nil
	case "classify-hot-discharge":
		err := fmt.Errorf("compressor protection: %w", &compressor.HotDischargeError{Temperature: 485})
		return map[string]any{"classification": compressor.Classify(err), "error": err.Error()}, nil
	case "classify-other":
		err := fmt.Errorf("compressor operation: %w", &compressor.OtherError{Message: detail})
		return map[string]any{"classification": compressor.Classify(err), "error": err.Error()}, nil
	case "converter-react":
		return map[string]any{"summary": p.converter.Summary()}, p.converter.React(ctx)
	case "converter-heat":
		return map[string]any{"valves": len(p.valves)}, converter.RunHeating(ctx, p.valves, -1)
	case "shutdown":
		optional := func(context.Context) error { return nil }
		quench := func(shutdownCtx context.Context) error { return p.cooling.Establish(shutdownCtx) }
		if err := converter.RunShutdown(ctx, optional, quench); err != nil {
			return nil, err
		}
		return map[string]any{"status": p.cooling.Status()}, StopComponents(ctx, optional, quench)
	case "separator-drain":
		return map[string]any{"timeout": separator.Timeout().String()}, p.separator.AuthorizeDrain(ctx)
	case "purge-apply":
		plan := []purge.ValveCommand{{Name: "purge-1", Position: 20}, {Name: "purge-2", Position: 25}}
		result := p.purgeDispatcher.Apply(plan)
		p.purgePlans.Commit(result.Snapshot)
		p.lastPurgeAudit = purge.NewAudit(result)
		return result, result.Err
	case "evidence-wait":
		current := make(chan string, 1)
		current <- "current-cycle"
		if err := (cooling.AnalyzerBarrier{}).Wait(ctx, []<-chan string{current}); err != nil {
			return nil, err
		}
		current <- "current-cycle"
		return map[string]any{"complete": true}, (purge.EvidenceBarrier{}).Wait(ctx, []<-chan string{current, nil})
	case "protect":
		if err := cooling.Integrate(ctx, p.cooling, p.telemetry); err != nil {
			return nil, err
		}
		equipment := []model.Equipment{feed.PermitEquipment(p.permit), compressor.EquipmentState(p.session), cooling.CondenserState(p.telemetry)}
		decision := interlock.Evaluate(equipment)
		if !decision.Allowed {
			p.interlock.SetEmergency(decision.Reason)
		}
		_ = campaign.IntegrateProtection(p.cycle, equipment)
		return map[string]any{"decision": decision, "cycle": interlock.Protect(p.cycle.Snapshot(), decision)}, nil
	case "recover":
		p.session.ClearTrip()
		equipment := model.RecoverEquipment(compressor.EquipmentState(p.session), time.Now().UTC())
		equipment = cooling.RecoverCooling(equipment)
		interlock.Recover(p.interlock, true)
		cycle := Restore(campaign.Recover(p.cycle.Snapshot()))
		return map[string]any{"equipment": equipment, "cycle": cycle}, nil
	case "disable-feed":
		p.permit.Disable()
		return feed.PermitEquipment(p.permit), nil
	case "journal-append":
		incident := journal.IncidentFor(detail, "operator")
		if err := p.journal.Append(incident); err != nil {
			return nil, err
		}
		p.metrics.CountIncident()
		return incident, nil
	case "codec-roundtrip":
		payload, err := model.EncodeCycle(p.startup.NewCycle())
		if err != nil {
			return nil, err
		}
		return model.DecodeCycle(payload)
	default:
		return nil, fmt.Errorf("unknown plant action %q", action)
	}
}

func (p *Plant) Status() map[string]any {
	p.mu.Lock()
	defer p.mu.Unlock()
	pressure, level := p.tank.Snapshot()
	emergency, reason := p.interlock.Snapshot()
	reading := p.analyzerState.Snapshot()
	product := separator.Separate(reading)
	incidents := p.journal.Entries()
	if replayed, err := journal.Replay(p.journalPath); err == nil {
		incidents = replayed
	}
	valves := make([]bool, len(p.valves))
	for index, valve := range p.valves {
		valves[index] = valve.IsOpen()
	}
	return map[string]any{
		"config_valid":      p.config.Valid(),
		"feed_gate":         p.quality.Describe(),
		"feed_pressure":     p.permit.Pressure(),
		"compressor":        compressor.EquipmentState(p.session),
		"cooling":           p.cooling.Status(),
		"separator":         map[string]any{"pressure": pressure, "level": level, "product_ready": separator.ProductReady(product)},
		"analysis":          reading,
		"converter_beds":    p.beds.Snapshot(),
		"converter_summary": p.converter.Summary(),
		"bypass_valves":     valves,
		"purge_plan":        p.purgePlans.Snapshot(),
		"purge_audit":       p.lastPurgeAudit,
		"interlock":         map[string]any{"emergency": emergency, "reason": reason},
		"incidents":         incidents,
		"metrics":           p.metrics.Snapshot(),
	}
}
