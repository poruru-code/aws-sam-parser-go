# Schema Coverage

This library ships Go types generated from the official AWS SAM JSON schema,
plus a small set of CloudFormation extensions maintained in this repository.

## Coverage policy
- **Base**: AWS SAM JSON schema (as provided by AWS).
- **Extensions**: Only the custom fragments stored under
  `tools/schema-gen/extensions/`.

The extensions are intentionally small and focused. If you need additional
CloudFormation resources or properties, add an extension and regenerate the
schema in the `aws-sam-parser-go` repository.

## Extension inventory
The authoritative list of extensions is the directory:
`tools/schema-gen/extensions/`

Use the filenames in that directory to see exactly what is added.
See `docs/schema-gen.md` for extension authoring and generator details.
