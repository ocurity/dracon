# Developing

We provide a complete development environment using the [Please build system](https://please.build). A `./pleasew` wrapper script is available for contributers who do not have Please installed and for use in CI/CD systems.

## Getting Started

1. Set up the development environment.

    ```bash
    $ ./pleasew dev
    ```

2. Deploy supporting resources that Dracon uses.

    ```bash
    $ ./pleasew dev_deploy
    ```

3. Make your changes :).
4. Run formatters.

    ```bash
    $ make fmt
    ```

5. Run linters.

    ```bash
    $ make lint
    ```

#### Cleaning Up

1. Run the following to delete the K3D cluster:

    ```bash
    $ ./pleasew run //build/k8s/k3d:teardown
    ```

2. Run the following to remove all build artefacts:

    ```bash
    $ ./pleasew cleanup
    ```
