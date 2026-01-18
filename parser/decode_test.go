// Where: parser/decode_test.go
// What: Unit tests for YAML decoding and intrinsic normalization.
// Why: Ensure stable decoding for SAM templates.
package parser

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDecodeYAML_Advanced(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]any
	}{
		{
			name: "scalar tags",
			content: `
A: !Ref MyParam
B: !Sub "hello-${Name}"
C: !GetAtt Res.Attr
D: !ImportValue ExportName
E: !Condition MyCond
`,
			want: map[string]any{
				"A": map[string]any{"Ref": "MyParam"},
				"B": map[string]any{"Fn::Sub": "hello-${Name}"},
				"C": map[string]any{"Fn::GetAtt": "Res.Attr"},
				"D": map[string]any{"Fn::ImportValue": "ExportName"},
				"E": map[string]any{"Condition": "MyCond"},
			},
		},
		{
			name: "sequence tags",
			content: `
A: !Join [ ":", [ a, b ] ]
B: !Equals [ a, b ]
C: !And [ true, false ]
D: !Or [ true, false ]
E: !Not [ true ]
F: !Select [ 0, [ a, b ] ]
G: !Split [ ":", "a:b" ]
H: !If [ Cond, yes, no ]
I: !Sub [ "hello-${Name}", { Name: world } ]
J: !GetAtt [ Res, Attr ]
`,
			want: map[string]any{
				"A": map[string]any{"Fn::Join": []any{":", []any{"a", "b"}}},
				"B": map[string]any{"Fn::Equals": []any{"a", "b"}},
				"C": map[string]any{"Fn::And": []any{true, false}},
				"D": map[string]any{"Fn::Or": []any{true, false}},
				"E": map[string]any{"Fn::Not": []any{true}},
				"F": map[string]any{"Fn::Select": []any{0, []any{"a", "b"}}},
				"G": map[string]any{"Fn::Split": []any{":", "a:b"}},
				"H": map[string]any{"Fn::If": []any{"Cond", "yes", "no"}},
				"I": map[string]any{"Fn::Sub": []any{"hello-${Name}", map[string]any{"Name": "world"}}},
				"J": map[string]any{"Fn::GetAtt": []any{"Res", "Attr"}},
			},
		},
		{
			name: "mapping tags",
			content: `
A: !Sub
  Name: world
  Template: hello-${Name}
`,
			want: map[string]any{
				"A": map[string]any{
					"Fn::Sub": map[string]any{
						"Name":     "world",
						"Template": "hello-${Name}",
					},
				},
			},
		},
		{
			name: "basic types",
			content: `
Int: 123
Float: 1.23
Bool: true
N: null
`,
			want: map[string]any{
				"Int":   123,
				"Float": 1.23,
				"Bool":  true,
				"N":     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeYAML(tt.content)
			if err != nil {
				t.Fatalf("DecodeYAML() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeYAML_NonStringKey(t *testing.T) {
	content := `
123: value
`
	got, err := DecodeYAML(content)
	if err != nil {
		t.Fatalf("DecodeYAML() error = %v", err)
	}
	want := map[string]any{"123": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DecodeYAML() = %v, want %v", got, want)
	}
}

func TestDecodeYAML_Errors(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"invalid yaml", "[:"},
		{"empty yaml", ""},
		{"sequence root", "- a\n- b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeYAML(tt.content)
			if err == nil {
				t.Errorf("DecodeYAML() expected error for %s", tt.name)
			}
		})
	}
}

func TestDecodeNode_DocumentNode(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: "!!str", Value: "A"},
					{Kind: yaml.ScalarNode, Tag: "!!str", Value: "B"},
				},
			},
		},
	}

	got := decodeNode(node)
	want := map[string]any{"A": "B"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("decodeNode() = %v, want %v", got, want)
	}
}

func TestDecodeScalar_Fallbacks(t *testing.T) {
	tests := []struct {
		name string
		node *yaml.Node
		want any
	}{
		{
			name: "invalid int falls back to string",
			node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "nope"},
			want: "nope",
		},
		{
			name: "invalid float falls back to string",
			node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: "nope"},
			want: "nope",
		},
		{
			name: "invalid bool falls back to string",
			node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "maybe"},
			want: "maybe",
		},
		{
			name: "unknown tag returns raw value",
			node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!Unknown", Value: "value"},
			want: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeScalar(tt.node)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeScalar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsString(t *testing.T) {
	if got := asString(nil); got != "" {
		t.Errorf("asString(nil) = %q, want empty string", got)
	}
	if got := asString(123); got != "123" {
		t.Errorf("asString(123) = %q, want %q", got, "123")
	}
}
