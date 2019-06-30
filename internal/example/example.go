package example

import (
	"context"

	"github.com/ben0x539/viewpoint"
	"github.com/ben0x539/viewpoint/observables/log"
	"github.com/ben0x539/viewpoint/observables/metrics"
)

var (
	Pkg               = viewpoint.MakePackage("github.com/ben0x539/viewpoint/internal/example")
	Log               = log.MakeLog(Pkg, "")
	CounterThings     = metrics.MakeCounter(Pkg, "things")
	GaugeThing        = metrics.MakeGauge(Pkg, "thing")
	DistributionThing = metrics.MakeDistribution(Pkg, "things")
)

func DoSomething(ctx context.Context) {
	Log.Debug(ctx, "unimportant things happened")
	Log.Info(ctx, "informative things happened")
	CounterThings.Count(ctx, 1)
	CounterThings.Count(ctx, 1)
	CounterThings.Count(ctx, 1)
	GaugeThing.Set(ctx, 41)
	GaugeThing.Set(ctx, 42)
	DistributionThing.Add(ctx, 1)
	DistributionThing.Add(ctx, 2)
	DistributionThing.Add(ctx, 3)
}
