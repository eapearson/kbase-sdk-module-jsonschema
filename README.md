# READ ME

This is a little server for JSON-SCHEMA.

## Hosting jsonschemas at KBase
Let’s make a jsonschema server, so. We can share our schemas.

## Features & Aspirations

- version schemas for safe schema evolution, adopting [SchemaVer](https://docs.snowplowanalytics.com/docs/pipeline-components-and-applications/iglu/common-architecture/schemaver/)
- access schemas by  specific version, or:
  - most recent `ADDITION` for a given `MODEL.REVISION`
  - most recent `REVISION.ADDITION` for a given `MODEL`
  - most recent `MODEL.REVISION.ADDITION`
- working, http-based refs
- support KBase schema extensions; e.g. ontology, taxonomy references, ui configuration
- provide (in this library, or others), jsonschema extensions to support the above.
- support for all KBase core types, e.g. KIDL boolean, workspace types

## Design
 
Schema urls look like:

`https://kbase.us/services/jsonschema/schema/<some domain path>/<schema>.<version>.json`

Where
- `services/jsonoschema` would be the core service path  (or if we had better dynamic service proxying, `/dynserv/jsonschema/`)
  `schema/` is the path within the service to invoke the jsonschema feature (there are other root paths, like `about`, `status`) 
- `<some domain path>` is some path, with one or more elements. E.g. foo/bar. The actual domain paths utilized should reflect the overall KBase namespacing; e.g. service typings at `/services/myservice`, KIDL core types  at  `kidl` , workspace types at `/workspace/types`.
- `<version>` is a [SchemaVer](https://docs.snowplowanalytics.com/docs/pipeline-components-and-applications/iglu/common-architecture/schemaver/) string; e.g. `1-2-3`
- `<schema>` is the name of the schema, and would correspond to type name, but in all lower case. 
  
## URLs

All path elements above should be in lower case.
Examples:
  `https://kbase.us/services/jsonschema/schemas/sample/fields/material.1-2-3.json`

Plural vs. singular? 

I'm thinking that path elements for which indicate collections are plural (e.g. "schemas") but those which indicate roughly a domain would be singular (e.g. `sample`)

Some alternative urls:
- `https://kbase.us/services/jsonchema/schemas/sample/fields/material.json` returns the most recent schema file
- `https://kbase.us/services/jsonschema/schemas/sample/fields/material.1.json` returns the most recent schema file with `1` as the MODEL.
- `https://kbase.us/services/jsonschema/schemas/sample/fields/material.1-2json` returns the most recent schema file with `1` as the MODEL and `2` as the REVISION.
- `https://kbase.us/services/schemas/sample/fields/material.1-2-3.json` returns the most exact schema file with `1` as the MODEL, `2` as the REVISION, and '3' as the ADDITION.

Good resource for versioning:
[https://snowplowanalytics.com/blog/2014/05/13/introducing-schemaver-for-semantic-versioning-of-schemas/]


Nah, let's just do JSON-RPC! We don't have to reinvent an rpc call in rest.

No, but we also need REST because regular GET requests is what feeds the beast (of $ref).

actual path
/schema/PATH/NAME.VERSION.json

start 

export SCHEMA_ROOT=`pwd`/schemas
go run server.go

try http://localhost:8080/sample/fields/material.json

## Development

### Build Executable

```shell
make compile
```

```shell
go build server.go
```

### Run server

```shell
export PORT=8080
go run server.go
```

or 

```shell
export PORT=8080
./server
```

### Build Development Image

```shell
make dev-image
```

### Run Development Image

```shell
export PORT=5000
make run-dev-image
```