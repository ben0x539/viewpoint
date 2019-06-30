package metrics

import (
	"context"

	"github.com/ben0x539/viewpoint"
)

type Counter struct {
	Label   string
	Package *viewpoint.Package
}

type Gauge struct {
	Label   string
	Package *viewpoint.Package
}

type Distribution struct {
	Label   string
	Package *viewpoint.Package
}

type CounterObserver interface {
	Inc(ctx context.Context, counter *Counter, count int64)
}

type GaugeObserver interface {
	Set(ctx context.Context, gauge *Gauge, value int64)
}

type DistributionObserver interface {
	Add(ctx context.Context, distribution *Distribution, value int64)
}

func MakeCounter(pkg *viewpoint.Package, label string) *Counter {
	return &Counter{
		Label:   label,
		Package: pkg,
	}
}

func MakeGauge(pkg *viewpoint.Package, label string) *Gauge {
	return &Gauge{
		Label:   label,
		Package: pkg,
	}
}

func MakeDistribution(pkg *viewpoint.Package, label string) *Distribution {
	return &Distribution{
		Label:   label,
		Package: pkg,
	}
}

func (c *Counter) GetParent() viewpoint.Observable {
	return c.Package
}

func (g *Gauge) GetParent() viewpoint.Observable {
	return g.Package
}

func (d *Distribution) GetParent() viewpoint.Observable {
	return d.Package
}

func (c *Counter) Count(ctx context.Context, count int64) {
	viewpoint.WithObserver(ctx, c, func(observer viewpoint.Observer) bool {
		counterObserver, ok := observer.(CounterObserver)
		if !ok {
			return false
		}

		counterObserver.Inc(ctx, c, count)
		return true
	})
}

func (g *Gauge) Set(ctx context.Context, value int64) {
	viewpoint.WithObserver(ctx, g, func(observer viewpoint.Observer) bool {
		gaugeObserver, ok := observer.(GaugeObserver)
		if !ok {
			return false
		}

		gaugeObserver.Set(ctx, g, value)
		return true
	})
}

func (d *Distribution) Add(ctx context.Context, value int64) {
	viewpoint.WithObserver(ctx, d, func(observer viewpoint.Observer) bool {
		distributionObserver, ok := observer.(DistributionObserver)
		if !ok {
			return false
		}

		distributionObserver.Add(ctx, d, value)
		return true
	})
}
