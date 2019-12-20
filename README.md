# Money

A project to keep track of the money I spent. All I need is an excel sheet, but
I figured I'd try to build something myself so that I could learn about
encryption and security.

## Contributing

### Conventions

```go
// yes
t.Errorf("should have ErrSomething, got %s", err)

// no
t.Errorf("wrong error, got %s, shoud have ErrSomething", err)
```

Tagging errors:

```go
var ErrTag = errors.New("some tag")

t.Errorf("my specific error %w", ErrTag)
```

Never call `r.ParseForm()`, but `r.ParseMultipartForm()`

## To think about

How should internal errors be handled? How can you design error handling to be
future proof (for eg. `MyFunc() error` never return an `ErrInternal` error in
the current version. What should my error handling look like? Do I just pass up
all the errors I don't know anything about? Should they be tagged then? Should
it even be legal to pass up an unknown error up the stack without tagging?)

A caller shouldn't wrap an error it doesn't know anything about

## Optimisation

This app is absolutely horrible when it comes to perfs. Probably the best thing
to do is to run it under a profiler to see what's the worse. In the mean time,
here are some ideas:

`users.list` should be a sorted list of independent JSON objects (not one big
lists), so that we can stop reading into memory as soon as we found what we
wanted...
