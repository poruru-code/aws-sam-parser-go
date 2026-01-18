// Where: parser/decoder_test.go
// What: Tests for Decode and DecodeToModel helpers.
// Why: Ensure decoded output matches expected Go types.
package parser

import (
	"testing"
)

type testModel struct {
	Name string `json:"Name"`
}

type countModel struct {
	Count int `json:"Count"`
}

func TestDecode_WithResolver(t *testing.T) {
	input := map[string]any{
		"Name": map[string]any{"Ref": "Env"},
	}

	resolver := refResolver{values: map[string]string{"Env": "dev"}}
	var out testModel
	if err := Decode(input, &out, &DecodeOptions{Resolver: resolver}); err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if out.Name != "dev" {
		t.Fatalf("unexpected Name: %s", out.Name)
	}
}

func TestDecode_WeaklyTyped(t *testing.T) {
	input := map[string]any{
		"Count": "123",
	}
	var out countModel
	if err := Decode(input, &out, &DecodeOptions{WeaklyTyped: true}); err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if out.Count != 123 {
		t.Fatalf("unexpected Count: %d", out.Count)
	}
}

func TestDecode_WeaklyTypedFalse(t *testing.T) {
	input := map[string]any{
		"Count": "123",
	}
	var out countModel
	if err := Decode(input, &out, &DecodeOptions{WeaklyTyped: false}); err == nil {
		t.Fatalf("expected Decode to fail for WeaklyTyped=false")
	}
}

func TestDecodeToModel_Minimal(t *testing.T) {
	content := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Hello:
    Type: AWS::Serverless::Function
    Properties:
      FunctionName: hello
`
	raw, err := DecodeYAML(content)
	if err != nil {
		t.Fatalf("DecodeYAML error: %v", err)
	}
	model, err := DecodeToModel(raw, nil)
	if err != nil {
		t.Fatalf("DecodeToModel error: %v", err)
	}
	if _, ok := model.Resources["Hello"]; !ok {
		t.Fatalf("expected Hello resource in model")
	}
}
