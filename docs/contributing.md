# Developing

Contributions to this project are more than welcome!

## Getting Started

0. The project is based in Go so you need to have the go binary installed.
   For linting Markdown files (required for component documentation) we use
   `remark` which is an `npm` package, so you will need that installed too.
   Make sure that you have a Go version equal or newer than the one listed in
   the `go.mod` file.
1. Make your changes :).
2. Run the formatters, linters and tests.

```bash
make fmt
```

4. Run linters.

```bash
make lint
```

5. Run tests.

```bash
make test
```

6. Make sure you have updated all manifests in case you are changing anything in
   the tools that generate templates.
7. All commits must be signed off and must include a ticket number along with a
   clear descriptionof what is being changed.
