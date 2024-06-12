# Building Protobufs

We use [buf.build](https://buf.build/) to simplify compiling the various `.proto` files.
You can install and compile the protos via:

```bash
make install-lint-tools
buf generate
```

This will update the generated protos in `api/proto/v1`.

## Adding new Protos

1. Add the proto to `api/proto/v1`
2. Append that file to the list of inputs in `buf.gen.yaml`
