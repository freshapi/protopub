# Protopub

Protopub allows you to publish your protobuf definitions into OCI complaint container registry
(such as docker-registry). This opens new possibilities for working with protoreflect - especially
dynamic interactive documentation and dynamic request routing.

## How to use

We have a simple tutorial [here](docs/tutorial.md), but here are some usage examples:

```bas
$ protopub build my-descriptor-set.bin                                        # build .proto files from current working directory into single descriptor set file using `protoc`
$ protopub login docker.io                                                    # login into registry
$ protopub push docker.io/freshapi/example:latest ./my-descriptor-set.bin     # push descriptor into registry
$ protopub inspect docker.io/freshapi/example:latest                          # get info about image
$ protopub pull docker.io/freshapi/example:latest ./pulled-descriptor-set.bin # pull from registry
```

That's it!

## Why

This project is heavily inspired by three factors:
1. A desire to build dynamic gRPC proxy (`grpc-router`)
2. A need for human-readable gRPC documentation in running system
3. [Buf](https://github.com/bufbuild/buf) project idea about protobuf schema registry. Technically this could be one of
the implementations of it, but surely we'd like to have custom backend which can perform certain validations
(e.g. backwards compatibility checks). Maybe be protopub could support Buf's image format in the future.

## What's next

Currently, there's no publicly available software which uses OCI registry to fetch protobuf schema because I am 
working on it right now. As soon as this project become available, this section will be updated. If you want
to use this mechanism in your internal project, feel free to import `pkg/protopub`.
