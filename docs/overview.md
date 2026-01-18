# Overview

`aws-sam-parser-go` is a small, composable library for working with AWS SAM
templates and CloudFormation-style YAML in Go.

It focuses on three responsibilities:

1. **Decode YAML into a neutral map**
   - Intrinsic tags (e.g., `!Ref`, `!Sub`) are normalized into
     canonical `Fn::` maps.
2. **Resolve values (optional)**
   - A resolver interface lets callers choose how to resolve intrinsics.
3. **Decode into generated types**
   - The `schema` package is generated from the AWS SAM JSON schema.

The core library intentionally avoids embedding project-specific assumptions.

## Packages
- `parser`: DecodeYAML, ResolveAll, Decode, and resolver interfaces.
- `schema`: Generated SAM schema types.
- `tools/schema-gen`: Schema merge and codegen tooling.

For precise behavior and limitations, see `docs/behavior.md`.
For schema coverage and extensions, see `docs/schema-coverage.md`.
