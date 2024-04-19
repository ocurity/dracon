# Developing

Contributions to this project are more than welcome!

## Getting Started

0. The project is based in Go so you need to have the go binary installed.
   For linting Markdown files (required for component documentation) we use
   `markdownlint-cli2` which is an `npm` package
1. Use the [Getting Started](../../docs/getting-started/installation.md) guide
   to setup your development environment and Tekton.
2. Make your changes :).
3. Run the formatters.

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
