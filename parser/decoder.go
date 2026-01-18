// Where: parser/decoder.go
// What: Decode helpers for mapping maps into Go structs.
// Why: Provide a consistent decoder for schema types.
package parser

import (
	"github.com/mitchellh/mapstructure"
	"github.com/poruru-code/aws-sam-parser-go/schema"
)

// DecodeOptions controls optional resolver behavior when decoding.
type DecodeOptions struct {
	Resolver        Resolver
	ResolverContext *Context
	WeaklyTyped     bool
}

// Decode maps input into output, optionally resolving values first.
func Decode(input any, output any, opts *DecodeOptions) error {
	resolved := input
	if opts != nil && opts.Resolver != nil {
		value, err := ResolveAll(opts.ResolverContext, input, opts.Resolver)
		if err != nil {
			return err
		}
		resolved = value
	}

	weaklyTyped := true
	if opts != nil {
		weaklyTyped = opts.WeaklyTyped
	}

	config := &mapstructure.DecoderConfig{
		Result:           output,
		TagName:          "json",
		WeaklyTypedInput: weaklyTyped,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(resolved)
}

// DecodeToModel decodes input into the generated SAM schema model.
func DecodeToModel(input any, opts *DecodeOptions) (schema.SamModel, error) {
	var model schema.SamModel
	if err := Decode(input, &model, opts); err != nil {
		return schema.SamModel{}, err
	}
	return model, nil
}
