// Where: examples/custom_resolver/main.go
// What: Custom Resolver implementation for Ref and Fn::If.
// Why: Demonstrate how callers can plug in their own intrinsic logic.
package main

import (
	"fmt"
	"log"

	"github.com/poruru-code/aws-sam-parser-go/parser"
)

type simpleResolver struct {
	refs       map[string]any
	conditions map[string]bool
}

func (r simpleResolver) Resolve(_ *parser.Context, value any) (any, bool, error) {
	m, ok := value.(map[string]any)
	if !ok || len(m) != 1 {
		return value, false, nil
	}

	if ref, ok := m["Ref"]; ok {
		name := fmt.Sprint(ref)
		if resolved, ok := r.refs[name]; ok {
			return resolved, true, nil
		}
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
	condName := fmt.Sprint(args[0])
	if r.conditions[condName] {
		return args[1], true, nil
	}
	return args[2], true, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("run error: %v", err)
	}
}

func run() error {
	input := map[string]any{
		"Name": map[string]any{"Ref": "Stage"},
		"Mode": map[string]any{"Fn::If": []any{"IsProd", "prod", "dev"}},
	}

	resolver := simpleResolver{
		refs:       map[string]any{"Stage": "staging"},
		conditions: map[string]bool{"IsProd": false},
	}

	resolved, err := parser.ResolveAll(&parser.Context{MaxDepth: 5}, input, resolver)
	if err != nil {
		return fmt.Errorf("ResolveAll error: %w", err)
	}

	fmt.Printf("resolved=%#v\n", resolved)
	return nil
}
