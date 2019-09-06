# Fhir json schema

## Example

Remove  property "_*" in json schema

```go

input_file := "PATH/fhir.schema.json"

output_file := "PATH/new.fhir.schema.json"

fhir_json_schema.RemoveLowerDash(input_file, output_file)

```

Generate Fhir Struct Map

```go

input_file := "PATH/fhir.schema.json"

output_file_struct_schema := "PATH/fhir-schema-struct-map.go"

fhir_json_schema.GenerateStructMap( input_file, output_file_struct_schema, "schema", "FhirStructMap")

```
