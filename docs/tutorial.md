# Protopub tutorial

Protopub is designed to publish your protobuf information to the OCI complaint
container registry (such as Docker Registry).

To start using `protopub` you will have to have container registry where you will push your
images.

## Step 1: clone repository

Clone protopub repository to gain access to the `examples/helloworld` directory:

```bash
$ git clone https://github.com/freshapi/protopub
$ cd protopub/examples
```

## Step 2: generate protobuf descriptor set

```bash
$ protopub build descriptor-set.bin
```

This will produce `descriptor-set.bin` file which we'll use to publish to the registry.

## Step 3: publish the descriptor set

You have to have a place where to publish your image. We will use docker.io/freshapi/example:latest as an example.
```bash
$ protopub login docker.io # you can skip this step if you have logged in with docker already
$ protopub push docker.io/freshapi/example:latest ./descriptor-set.bin
```

## Step 4: download pushed descriptor set

You can now pull back descriptor set from the registry:
```bash
$ protopub pull docker.io/freshapi/example:latest ./descriptor-set-pulled.bin
```
