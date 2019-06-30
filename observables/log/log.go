package log

import (
	"context"
	"os"

	"github.com/ben0x539/viewpoint"
)

type Log struct {
	Name    string
	Label   string
	Package *viewpoint.Package
}

type LogObserver interface {
	Log(ctx context.Context, log *Log, level Level, msg string)
}
type Level int

const (
	LevelDebug Level = iota
	LevelInfo  Level = iota
	LevelWarn  Level = iota
	LevelError Level = iota
	LevelFatal Level = iota
)

func MakeLog(pkg *viewpoint.Package, name string) *Log {
	var label string
	if pkg != nil && name != "" {
		label = pkg.Name + "/" + name
	} else if pkg != nil {
		label = pkg.Name
	} else {
		label = name
	}

	return &Log{
		Name:    name,
		Label:   label,
		Package: pkg,
	}
}

func (l *Log) GetParent() viewpoint.Observable {
	if l.Package == nil {
		return nil
	}

	return l.Package
}

func (l *Log) Log(ctx context.Context, level Level, msg string) {
	viewpoint.WithObserver(ctx, l, func(observer viewpoint.Observer) bool {
		logObserver, ok := observer.(LogObserver)
		if !ok {
			return false
		}

		logObserver.Log(ctx, l, level, msg)
		return true
	})
}

func (l *Log) Info(ctx context.Context, msg string) {
	l.Log(ctx, LevelInfo, msg)
}

func (l *Log) Debug(ctx context.Context, msg string) {
	l.Log(ctx, LevelDebug, msg)
}

func (l *Log) Warn(ctx context.Context, msg string) {
	l.Log(ctx, LevelWarn, msg)
}

func (l *Log) Error(ctx context.Context, msg string) {
	l.Log(ctx, LevelError, msg)
}

func (l *Log) Fatal(ctx context.Context, msg string) {
	l.Log(ctx, LevelFatal, msg)
	os.Exit(1)
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		panic("weird log level")
	}
}
