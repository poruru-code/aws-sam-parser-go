# aws-sam-parser-go

A minimal Go helper library for parsing AWS SAM templates and CloudFormation-style YAML.

It provides:
- YAML decoding that normalizes intrinsic tags into canonical `Fn::` maps.
- A resolver interface and recursive walker (caller-defined behavior).
- Generated Go structs for the AWS SAM JSON schema.

Intrinsic resolution only happens when a resolver is provided via
`parser.ResolveAll` or `parser.Decode`/`DecodeOptions`. Without a resolver the
walker simply emits the normalized `Fn::` maps. Pass a `Context` with
`MaxDepth` (e.g., `MaxDepth: 20`) to prevent infinite recursion when your
resolver can revisit the same tree nodes. See `docs/resolver.md` for examples.

## Packages
- `parser`: YAML decode, resolver interface, tree walker, and decode helpers.
- `schema`: Generated Go types for the SAM schema.
- `tools/schema-gen`: Schema merge + codegen helpers.

## Installation
```bash
go get github.com/poruru-code/aws-sam-parser-go
```

## Quick Start
```go
raw, err := parser.DecodeYAML(content)
if err != nil {
	return err
}

// Optional: resolve intrinsics using your own resolver.
resolved, err := parser.ResolveAll(&parser.Context{MaxDepth: 20}, raw, myResolver)
if err != nil {
	return err
}

var model schema.SamModel
if err := parser.Decode(resolved, &model, nil); err != nil {
	return err
}
```

Alternatively, pass a resolver directly to Decode:
```go
var model schema.SamModel
if err := parser.Decode(raw, &model, &parser.DecodeOptions{
	Resolver:        myResolver,
	ResolverContext: &parser.Context{MaxDepth: 20},
}); err != nil {
	return err
}
```

## Examples
- `examples/decode_resolve_model`: DecodeYAML -> ResolveAll -> DecodeToModel pipeline.
- `examples/custom_resolver`: Custom Resolver for `Ref` and `Fn::If`.

## Behavior
Detailed behavior, supported tags, and limitations are documented in
`docs/behavior.md`.

Highlights:
- Supported intrinsic tags and normalization rules
- DecodeYAML error cases and scalar conversion
- ResolveAll traversal and depth limits

## Schema Coverage
Coverage and extension policy are documented in `docs/schema-coverage.md`.

If a CloudFormation resource is missing (e.g., `AWS::SNS::Topic`), add an
extension and regenerate the schema. See `docs/schema-gen.md` for the
step-by-step flow.

## Resolver Design
The library does not hardcode intrinsic resolution rules.
Instead, you provide a `Resolver` implementation and pass it via `ResolveAll`
or `DecodeOptions`:

```go
type Resolver interface {
	Resolve(ctx *Context, value any) (out any, handled bool, err error)
}
```

This keeps the core neutral and allows you to plug in:
- Strict CloudFormation-compatible resolution
- Project-specific behaviors
- Partial resolution (e.g., only `Ref` or `Sub`)

See `docs/resolver.md` for a small example.

## Schema Generation
Generated schema lives in `schema/sam_generated.go`.

```bash
cd tools/schema-gen
python3 generate.py
```

The generator merges the base SAM schema with custom extensions under
`tools/schema-gen/extensions/` (currently covering only `AWS::S3::Bucket` and
`AWS::DynamoDB::Table` fragments) and writes to `schema/sam_generated.go`.

### Generator version
Pin the generator to a known release for deterministic output:

```bash
go install github.com/elastic/go-json-schema-generate/cmd/schema-generate@v0.0.0-20220519132038-c708d18d6ca2
```

Update this pinned version when you intentionally bump the generator.

## Development
This repo uses `mise` for toolchains and `lefthook` for git hooks.

```bash
mise install
mise run setup
mise run check
```

## License
Apache-2.0. See `LICENSE`.
