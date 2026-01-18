# Resolver

The resolver is an optional layer that converts intrinsic maps into concrete
values. The core does not implement any intrinsic rules by default, so
resolution only happens when you provide a resolver.

When calling `ResolveAll`, pass a `Context` with `MaxDepth` set to guard
against recursive definitions (ESB uses `MaxDepth=20`). Once `depth >
MaxDepth` the walker returns the raw node without invoking the resolver,
preventing unbounded recursion.

## Interface
```go
type Resolver interface {
	Resolve(ctx *Context, value any) (out any, handled bool, err error)
}
```

- `handled` should be true when the resolver produced a new value.
- If `handled` is false, the walker continues with child values as-is.

## Minimal Example
```go
type RefResolver struct {
	Values map[string]string
}

func (r RefResolver) Resolve(_ *Context, value any) (any, bool, error) {
	m, ok := value.(map[string]any)
	if !ok || len(m) != 1 {
		return value, false, nil
	}
	ref, ok := m["Ref"]
	if !ok {
		return value, false, nil
	}
	name := fmt.Sprint(ref)
	if v, ok := r.Values[name]; ok {
		return v, true, nil
	}
	return value, false, nil
}
```

You can chain resolvers by implementing a simple dispatcher that tries each
resolver in order.
