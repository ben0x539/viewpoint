package main

import (
	"context"
	"expvar"
	"log"

	"github.com/ben0x539/viewpoint"
	v_log "github.com/ben0x539/viewpoint/observables/log"
	v_metrics "github.com/ben0x539/viewpoint/observables/metrics"

	v_expvar "github.com/ben0x539/viewpoint/observers/expvar"
	v_stdlog "github.com/ben0x539/viewpoint/observers/stdlog"

	"github.com/ben0x539/viewpoint/internal/example"
)

type MyAppObserver struct {
	*v_stdlog.Logger
	*v_expvar.ExpvarMetrics
}

var _ v_metrics.CounterObserver = &MyAppObserver{}
var _ v_log.LogObserver = &MyAppObserver{}

func main() {
	expvarMetrics := &v_expvar.ExpvarMetrics{}
	expvarMetrics.Publish(example.Pkg.Name + ".stats")
	myAppObserver := &MyAppObserver{
		Logger:        v_stdlog.MakeLogger(v_log.LevelInfo),
		ExpvarMetrics: expvarMetrics,
	}
	ctx := viewpoint.Configure(context.Background(), viewpoint.Observers{
		example.Pkg: myAppObserver,
	}, v_stdlog.MakeLogger(v_log.LevelWarn))

	example.DoSomething(ctx)

	UnrelatedLog := v_log.MakeLog(nil, "weird other log")

	UnrelatedLog.Info(ctx, "unrelated info")
	UnrelatedLog.Warn(ctx, "unrelated warn")

	log.Printf("example pkg expvar contents: %v", expvar.Get("github.com/ben0x539/viewpoint/internal/example.stats"))
}
