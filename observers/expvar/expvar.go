package expvar

import (
	"context"
	"expvar"
	"sync"
	"sync/atomic"

	v_metrics "github.com/ben0x539/viewpoint/observables/metrics"
)

type ExpvarMetrics struct {
	counters      sync.Map
	gauges        sync.Map
	distributions sync.Map
}

type sliceWithMutex struct {
	values []int64
	mu     sync.Mutex
}

type expData struct {
	Counters      map[string]int64
	Gauges        map[string]int64
	Distributions map[string][]int64
}

var _ v_metrics.CounterObserver = &ExpvarMetrics{}

func (e *ExpvarMetrics) Inc(ctx context.Context, counter *v_metrics.Counter, count int64) {
	valueInterface, ok := e.counters.Load(counter.Label)
	if !ok {
		var value int64
		valueInterface, _ = e.counters.LoadOrStore(counter.Label, &value)
	}

	atomic.AddInt64(valueInterface.(*int64), count)
}

func (e *ExpvarMetrics) Set(ctx context.Context, gauge *v_metrics.Gauge, value int64) {
	valueInterface, ok := e.gauges.LoadOrStore(gauge.Label, &value)
	if !ok {
		return
	}

	atomic.StoreInt64(valueInterface.(*int64), value)
}

func (e *ExpvarMetrics) Add(ctx context.Context, distribution *v_metrics.Distribution, value int64) {
	valueInterface, ok := e.distributions.Load(distribution.Label)
	if !ok {
		var values sliceWithMutex
		valueInterface, _ = e.distributions.LoadOrStore(distribution.Label, &values)
	}

	values := valueInterface.(*sliceWithMutex)
	values.mu.Lock()
	defer values.mu.Unlock()
	values.values = append(values.values, value)
}

func (e *ExpvarMetrics) Publish(name string) {
	expvar.Publish(name, expvar.Func(func() interface{} {
		exp := expData{
			Counters:      map[string]int64{},
			Gauges:        map[string]int64{},
			Distributions: map[string][]int64{},
		}

		e.counters.Range(func(k, v interface{}) bool {
			exp.Counters[k.(string)] = atomic.LoadInt64(v.(*int64))
			return true
		})

		e.gauges.Range(func(k, v interface{}) bool {
			exp.Gauges[k.(string)] = atomic.LoadInt64(v.(*int64))
			return true
		})

		e.distributions.Range(func(k, v interface{}) bool {
			values := v.(*sliceWithMutex)
			values.mu.Lock()
			defer values.mu.Unlock()

			exp.Distributions[k.(string)] = append([]int64{}, values.values...)

			return true
		})

		return exp
	}))
}
