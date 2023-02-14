# Dracon Demos

Example of running the the [Demo Dracon Pipelines](examples/pipelines/golang-project) where you get to see the results in Kibana at the end.

## Prerequisites

- You have followed the [Getting Started with k3d Guide](/docs/getting-started/k3d.md).

---

## Tutorial

1. You can run a demo pipeline with:

   ```bash
   $ ./pleasew deploy //examples/pipelines/golang-project:pipeline
   $ ./pleasew deploy //examples/pipelines/golang-project/pipelinerun:pipelinerun
   ```

2. Wait for the pipeline to finish running by monitoring it in https://tekton.dracon.localhost:8443.

3. Once the pipelinerun has finished running you can view your results in Kibana: https://kibana.dracon.localhost:8443.
