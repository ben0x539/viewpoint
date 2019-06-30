# viewpoint

Sometimes I write go programs. I'd like my `main` package to configure logs and
metrics emitted by code in other packages. I don't like passing explicit
loggers and statsd clients all over the place. I also don't like mutable global
variables.

I do like passing `context.Context` objects all over the place, and I guess
static, read-only global variables are fine. Maybe I can come up with a clever
scheme to hide the loggers and everything in a `context.Context` keyed by
values defined in each package.

The idea is for `main` to say something like

```go
ctx := viewpoint.Configure(ctx, viewpoint.Observers{
  some_package.Pkg: logAndSendToStatsd,
  some_other_package.Pkg: onlyLogWarnings,
  some_other_package.SomeSpecificCounter: sendToStatsd,
}, ...)

go some_package.Run(ctx)
go some_other_package.Run(ctx)
```

and have everything work out like you'd expect from reading that.

This whole thing is extremely 0.x, I'm not really happy with most of the names,
the way I'm using interfaces kinda sucks for composition, and everything is
just minimally interesting strawman implementations of things you might want to
observe, so I'll either never touch this again or make reckless breaking
changes.

It's called `viewpoint` because a package can declare observables by itself,
but what that _means_ depends on just how the `main` package wants to look at
them.
