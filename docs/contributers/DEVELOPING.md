# Developing

We provide a complete development environment using the [Please build system](https://please.build). A `./pleasew` wrapper script is available for contributers who do not have Please installed and for use in CI/CD systems.

## Getting Started

1. Use the [Getting Started](../../docs/getting-started/installation.md) guide to setup your
development environment and Tekton.
2. Make your changes :).
3. Run the formatters.

    ```bash
    $ make fmt
    ```

4. Run linters.

    ```bash
    $ make lint
    ```

5. Run tests.

    ```bash
    $ make test
    ```

6. Make sure you have updated all manifests in case you are changing anything in the tools that
generate templates.

7. All commits must be signed off and must include a ticket number along with a clear description
of what is being changed.
