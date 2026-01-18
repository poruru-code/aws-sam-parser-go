// Where: parser/resolve_test.go
// What: Tests for resolver traversal behavior.
// Why: Ensure resolvers can be composed and applied recursively.
package parser

import (
	"fmt"
	"reflect"
	"testing"
)

type refResolver struct {
	values map[string]string
}

func (r refResolver) Resolve(_ *Context, value any) (any, bool, error) {
	m, ok := value.(map[string]any)
	if !ok || len(m) != 1 {
		return value, false, nil
	}
	ref, ok := m["Ref"]
	if !ok {
		return value, false, nil
	}
	name := fmt.Sprint(ref)
	resolved, ok := r.values[name]
	if !ok {
		return value, false, nil
	}
	return resolved, true, nil
}

type ifResolver struct {
	flag bool
}

func (r ifResolver) Resolve(_ *Context, value any) (any, bool, error) {
	m, ok := value.(map[string]any)
	if !ok || len(m) != 1 {
		return value, false, nil
	}
	raw, ok := m["Fn::If"]
	if !ok {
		return value, false, nil
	}
	args, ok := raw.([]any)
	if !ok || len(args) != 3 {
		return value, false, nil
	}
	if r.flag {
		return args[1], true, nil
	}
	return args[2], true, nil
}

type chainResolver []Resolver

func (c chainResolver) Resolve(ctx *Context, value any) (any, bool, error) {
	for _, resolver := range c {
		out, handled, err := resolver.Resolve(ctx, value)
		if err != nil {
			return nil, true, err
		}
		if handled {
			return out, true, nil
		}
	}
	return value, false, nil
}

func TestResolveAll_RefAndIf(t *testing.T) {
	input := map[string]any{
		"Value": map[string]any{
			"Fn::If": []any{
				"Cond",
				map[string]any{"Ref": "Env"},
				"fallback",
			},
		},
	}

	resolver := chainResolver{
		ifResolver{flag: true},
		refResolver{values: map[string]string{"Env": "dev"}},
	}

	resolved, err := ResolveAll(&Context{MaxDepth: 10}, input, resolver)
	if err != nil {
		t.Fatalf("ResolveAll error: %v", err)
	}

	want := map[string]any{"Value": "dev"}
	if !reflect.DeepEqual(resolved, want) {
		t.Fatalf("ResolveAll = %v, want %v", resolved, want)
	}
}

func TestResolveAll_MaxDepth(t *testing.T) {
	inner := map[string]any{"Ref": "Deep"}
	data := map[string]any{
		"outer": map[string]any{
			"inner": inner,
		},
	}

	resolver := refResolver{values: map[string]string{"Deep": "done"}}
	resolved, err := ResolveAll(&Context{MaxDepth: 1}, data, resolver)
	if err != nil {
		t.Fatalf("ResolveAll error: %v", err)
	}

	outer, ok := resolved.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", resolved)
	}
	innerMap, ok := outer["outer"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map, got %#v", outer["outer"])
	}
	if _, handled := innerMap["inner"].(map[string]any); !handled {
		t.Fatalf("expected inner map to remain unresolved, got %T", innerMap["inner"])
	}
}

func TestResolveAll_NilResolver(t *testing.T) {
	input := map[string]any{
		"Value": map[string]any{
			"Ref": "Env",
		},
	}

	resolved, err := ResolveAll(nil, input, nil)
	if err != nil {
		t.Fatalf("ResolveAll error: %v", err)
	}
	if !reflect.DeepEqual(resolved, input) {
		t.Fatalf("ResolveAll = %v, want %v", resolved, input)
	}
}

func TestResolveAll_NilContext(t *testing.T) {
	input := map[string]any{"Ref": "Env"}
	resolver := refResolver{values: map[string]string{"Env": "dev"}}

	resolved, err := ResolveAll(nil, input, resolver)
	if err != nil {
		t.Fatalf("ResolveAll error: %v", err)
	}
	if resolved != "dev" {
		t.Fatalf("ResolveAll = %v, want %v", resolved, "dev")
	}
}

func TestResolveAll_SliceTraversal(t *testing.T) {
	input := []any{
		map[string]any{"Ref": "Env"},
		"keep",
	}
	resolver := refResolver{values: map[string]string{"Env": "dev"}}

	resolved, err := ResolveAll(nil, input, resolver)
	if err != nil {
		t.Fatalf("ResolveAll error: %v", err)
	}
	want := []any{"dev", "keep"}
	if !reflect.DeepEqual(resolved, want) {
		t.Fatalf("ResolveAll = %v, want %v", resolved, want)
	}
}

func TestResolveAll_ResolverError(t *testing.T) {
	testErr := fmt.Errorf("boom")
	resolver := chainResolver{
		errResolver{err: testErr},
	}
	_, err := ResolveAll(&Context{MaxDepth: 10}, map[string]any{"Ref": "fail"}, resolver)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != testErr.Error() {
		t.Fatalf("unexpected error: %v", err)
	}
}

type errResolver struct {
	err error
}

func (r errResolver) Resolve(_ *Context, _ any) (any, bool, error) {
	return nil, true, r.err
}
