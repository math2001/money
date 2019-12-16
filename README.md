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

## To think about

How should internal errors be handled? How can you design error handling to be
future proof (for eg. `MyFunc() error` never return an `ErrInternal` error in
the current version. What should my error handling look like? Do I just pass up
all the errors I don't know anything about? Should they be tagged then? Should
it even be legal to pass up an unknown error up the stack without tagging?)

A caller shouldn't wrap an error it doesn't know anything about
