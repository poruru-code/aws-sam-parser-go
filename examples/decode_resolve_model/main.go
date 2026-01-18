// Where: examples/decode_resolve_model/main.go
// What: End-to-end example using DecodeYAML -> ResolveAll -> DecodeToModel.
// Why: Show the minimal pipeline for SAM template parsing.
package main

import (
	"fmt"
	"log"

	"github.com/poruru-code/aws-sam-parser-go/parser"
)

type refResolver struct {
	values map[string]string
}

func (r refResolver) Resolve(_ *parser.Context, value any) (any, bool, error) {
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

func main() {
	if err := run(); err != nil {
		log.Fatalf("run error: %v", err)
	}
}

func run() error {
	content := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  Hello:
    Type: AWS::Serverless::Function
    Properties:
      FunctionName: !Ref Env
`

	raw, err := parser.DecodeYAML(content)
	if err != nil {
		return fmt.Errorf("DecodeYAML error: %w", err)
	}

	resolver := refResolver{values: map[string]string{"Env": "dev"}}
	resolved, err := parser.ResolveAll(&parser.Context{MaxDepth: 10}, raw, resolver)
	if err != nil {
		return fmt.Errorf("ResolveAll error: %w", err)
	}

	model, err := parser.DecodeToModel(resolved, &parser.DecodeOptions{WeaklyTyped: true})
	if err != nil {
		return fmt.Errorf("DecodeToModel error: %w", err)
	}

	for name, resource := range model.Resources {
		if m, ok := resource.(map[string]any); ok {
			if typ, ok := m["Type"]; ok {
				fmt.Printf("resource=%s type=%v\n", name, typ)
			} else {
				fmt.Printf("resource=%s data=%v\n", name, m)
			}
		} else {
			fmt.Printf("resource=%s data=%T\n", name, resource)
		}
		break
	}

	return nil
}
