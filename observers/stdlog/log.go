package stdlog

import (
	"context"
	"log"

	v_log "github.com/ben0x539/viewpoint/observables/log"
)

type Logger struct {
	MinLevel v_log.Level
}

var _ v_log.LogObserver = &Logger{}

func (l *Logger) Log(ctx context.Context, log_ *v_log.Log, level v_log.Level, msg string) {
	if level < l.MinLevel {
		return
	}

	log.Printf("[%v] %v: %v", level, log_.Label, msg)
}

func MakeLogger(minLevel v_log.Level) *Logger {
	return &Logger{MinLevel: minLevel}
}
