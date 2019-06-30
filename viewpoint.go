package viewpoint

import (
	"context"
)

type Observable interface {
	GetParent() Observable
}

type Package struct {
	Name   string
	Parent *Package
}

func MakePackage(name string) *Package {
	return &Package{
		Name: name,
	}
}

func (p *Package) GetParent() Observable {
	if p.Parent == nil {
		return nil
	}

	return p.Parent
}

func (p *Package) SubPackage(name string) *Package {
	return &Package{
		Name:   name,
		Parent: p,
	}
}

type Observer interface{}

type Observers map[Observable]Observer

var contextConfigurationKey = struct{}{}

type ContextConfiguration struct {
	Observers Observers
	Fallback  Observer
	Parent    *ContextConfiguration
}

func Configure(ctx context.Context, observers Observers, fallback Observer) context.Context {
	superInterface := ctx.Value(contextConfigurationKey)
	super, _ := superInterface.(*ContextConfiguration)

	return context.WithValue(ctx, contextConfigurationKey, &ContextConfiguration{
		Observers: observers,
		Fallback:  fallback,
		Parent:    super,
	})
}

func WithObserver(ctx context.Context, observable Observable, f func(Observer) bool) {
	cfgInterface := ctx.Value(contextConfigurationKey)
	outerCfg, _ := cfgInterface.(*ContextConfiguration)

	for o := observable; o != nil; o = o.GetParent() {
		for cfg := outerCfg; cfg != nil; cfg = cfg.Parent {
			observer, ok := cfg.Observers[o]
			if !ok {
				continue
			}

			if f(observer) {
				return
			}
		}
	}

	for o := observable; o != nil; o = o.GetParent() {
		for cfg := outerCfg; cfg != nil; cfg = cfg.Parent {
			if cfg.Fallback == nil {
				continue
			}

			if f(cfg.Fallback) {
				return
			}
		}
	}
}
