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
