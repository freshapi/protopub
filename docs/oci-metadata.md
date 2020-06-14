# OCI metadata

## Config file

OCI image config file contains this JSON fields:

|Field|Value|Example|
|-----|-----|-------|
| freshapi.protopub.version | protopub schema version | "v1" 
| freshapi.protopub.files | An array containing all proto file names in this image | `["helloworld/helloworld.proto", "routeguide/route_guide.proto"]`

## Layer

Each OCI image layer represents single `FileDescriptor` object to allow partial loading of required descriptor set. Also,
each layer contains metadata in form of annotations.

### Annotations

|Annotation|Value|Example|
|----------|-----|-------|
| freshapi.protopub.filename | Name of proto file in descriptor set | helloworld/helloworld.proto |
| freshapi.protopub.syntax | File syntax | proto3 |
| freshapi.protopub.package | Package name of proto file | helloworld |
| freshapi.protopub.imports | Imported file names separated by comma | routeguide/route_guide.proto,helloworld/helloworld.proto |
| freshapi.protopub.messages | Message names separated by comma | HelloRequest,HelloReply |
| freshapi.protopub.enums | Enum names separated by comma | MyEnum1,MyEnum2 |
| freshapi.protopub.services | Service names separated by comma | Greeter |
| freshapi.protopub.sourcecodeinfo | True if descriptor contains SourceCodeInfo | true |
