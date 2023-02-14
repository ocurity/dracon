## Getting Started with Dracon with Please on [k3d](https://k3d.io)

k3d is a lightweight wrapper to run [k3s](https://k3s.io/) (Rancher Lab’s minimal Kubernetes distribution) in docker.

k3d makes it very easy to create single- and multi-node k3s clusters in docker, e.g. for local development on Kubernetes.

## Quick Guide

1. Run the helper script which creates and configures a development k3d cluster.

    ```bash
    $ ./pleasew dev
    ```

2. Deploy supporting resources that Dracon uses.

    ```bash
    $ ./pleasew dev_deploy
    ```

3. Dracon is now ready to use. You can see available endpoints/ingresses by
    running:

    ```bash
    $ kubectl get ingress --all-namespaces -o jsonpath='{range .items[*]}{.spec.rules[0].host}{"\n"}{end}' \
    | sed 's/\(dracon\.localhost$\)/\1\:8080/g' \
    | sed 's/^/http:\/\//g'
    ```
    Check out the [Running Demos Guide](/docs/getting-started/tutorials/running-demos.md)

### Running an Example Pipeline

1. Create the Pipeline in Tekton:

    ```bash
    $ ./pleasew run //examples/pipelines/golang-project:golang-project_deploy
    ```

2. Run the Pipeline in Tekton:

    ```bash
    $ ./pleasew run //examples/pipelines/golang-project:golang-project_pipelinerun_deploy
    ```

### Inspecting a Pipeline

We can view the progress of a pipeline in 2 ways:

- Visit and browse http://tekton.dracon.localhost:8080
- Use `kubectl get pipelineruns`
