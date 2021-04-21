# READ ME

This is a little server for JSON-SCHEMA.

# Hosting jsonschemas at KBase
Let’s make a jsonschema server, so. We can share our schemas.

## Design
It isn’t complicated, really. Just an HTTP server with a given directory structure.
`/schemas/<some domain path>/<version>/<schema>.json`
Where
- `<some domain path>` is some path, with one or more elements. E.g. foo/bar.
- `<version>` is a [SchemaVer] string; e.g. `1-2-3`
- `<schema>` is the name of the schema
  All path elements above should be in lower case.
  Examples:
  `https://kbase.us/services/schemas/sample/fields/1-2-3/material.json`

Some alternative “smart” urls:
- `https://kbase.us/services/schemas/sample/fields/material.json` returns the most recent schema file
- `https://kbase.us/services/schemas/sample/fields/1/material.json` returns the most recent schema file with `1` as the MODEL.
- `https://kbase.us/services/schemas/sample/fields/1-2/material.json` returns the most recent schema file with `1` as the MODEL and `2` as the REVISION.

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