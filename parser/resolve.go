// Where: parser/resolve.go
// What: Resolver interface and recursive walker.
// Why: Allow caller-defined intrinsic resolution without hardcoding rules.
package parser

// Resolver resolves values when possible and signals if it handled the input.
type Resolver interface {
	Resolve(ctx *Context, value any) (out any, handled bool, err error)
}

// Context controls resolver behavior during traversal.
type Context struct {
	MaxDepth int
}

// ResolveAll walks the value tree and applies the resolver when provided.
func ResolveAll(ctx *Context, value any, resolver Resolver) (any, error) {
	return resolveAll(ctx, value, resolver, 0)
}

func resolveAll(ctx *Context, value any, resolver Resolver, depth int) (any, error) {
	if ctx != nil && ctx.MaxDepth > 0 && depth > ctx.MaxDepth {
		return value, nil
	}

	if resolver != nil {
		out, handled, err := resolver.Resolve(ctx, value)
		if err != nil {
			return nil, err
		}
		if handled {
			return resolveAll(ctx, out, resolver, depth+1)
		}
	}

	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, v := range typed {
			resolved, err := resolveAll(ctx, v, resolver, depth+1)
			if err != nil {
				return nil, err
			}
			out[k] = resolved
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for i, v := range typed {
			resolved, err := resolveAll(ctx, v, resolver, depth+1)
			if err != nil {
				return nil, err
			}
			out[i] = resolved
		}
		return out, nil
	default:
		return value, nil
	}
}
