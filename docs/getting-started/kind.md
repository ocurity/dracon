## Getting Started with Dracon on [KinD](https://kind.sigs.k8s.io/)

KinD is is a tool for running local Kubernetes clusters using Docker container “nodes”.

## Quick Guide


1. Create KinD cluster named `dracon-demo`. For more info, see [official documentation](https://kind.sigs.k8s.io/docs/user/quick-start/#creating-a-cluster):

    ```bash
    $ kind create cluster --name dracon-demo
    ```

2. Install Tekton Pipelines. For more info, see [official documentation](https://tekton.dev/docs/installation/pipelines/#installing-tekton-pipelines-on-kubernetes).

    ```bash
    $ kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
    ```

3. **Optional**: Install Tekton Dashboard for a Web UI. For more info, see [official documentation](https://github.com/tektoncd/dashboard/blob/main/docs/install.md).

    ```bash
    $ curl -sL https://raw.githubusercontent.com/tektoncd/dashboard/main/scripts/release-installer | \
        bash -s -- install latest --read-write
    # Use `kubectl proxy` so you can access Kubernetes services on your local machine.
    $ kubectl proxy
    # Tekton Dashboard is now available at: http://localhost:8001/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/about
    ```

4. **Optional**: Install ECK and create an Elasticsearch + Kibana Dashboards. For more info, see [official documentation](https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-deploy-eck.html).

    ```bash
    # Create ECK CRDs.
    $ kubectl create -f https://download.elastic.co/downloads/eck/2.6.1/crds.yaml
    # Apply ECK operator resources.
    $ kubectl apply -f https://download.elastic.co/downloads/eck/2.6.1/operator.yaml
    # Create Elasticsearch.
    $ cat <<EOF | kubectl apply -f -
    apiVersion: elasticsearch.k8s.elastic.co/v1
    kind: Elasticsearch
    metadata:
    name: quickstart
    spec:
    version: 8.6.1
    nodeSets:
    - name: default
        count: 1
        config:
        node.store.allow_mmap: false
    EOF
    # Create Kibana.
    $ cat <<EOF | kubectl apply -f -
    apiVersion: kibana.k8s.elastic.co/v1
    kind: Kibana
    metadata:
    name: quickstart
    spec:
    version: 8.6.1
    count: 1
    elasticsearchRef:
        name: quickstart
    EOF
    ```

### Composing a Pipeline

We use [Kustomize Components](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/components.md) to create composable Dracon Pipelines.

1. Create the following simple Dracon Pipeline in your directory:

    ```yaml
    ---
    # ./kustomization.yaml
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization

    nameSuffix: -github-com-kubernetes-kubernetes
    namespace: default

    resources:
    - https://github.com/ocurity/dracon//components/base/

    components:
    - https://github.com/ocurity/dracon//components/sources/git/
    - https://github.com/ocurity/dracon//components/producers/aggregator/
    - https://github.com/ocurity/dracon//components/producers/golang-gosec/
    - https://github.com/ocurity/dracon//components/producers/golang-nancy/
    - https://github.com/ocurity/dracon//components/enrichers/aggregator/
    - https://github.com/ocurity/dracon//components/enrichers/deduplication/
    - https://github.com/ocurity/dracon//components/consumers/elasticsearch/
    ```

2. Run the following to create the Tekton Pipeline, Task, etc. resources on your cluster:

    ```bash
    $ kustomize build | kubectl apply -f -
    # Note: you can just run the below to see the generated Tekton Pipeline resources
    # $ kustomize build
    ```

3. Create the following Tekton PipelineRun file:

    ```yaml
    ---
    # pipelinerun.yaml
    # Run `kubectl create -f -` with this file.
    apiVersion: tekton.dev/v1beta1
    kind: PipelineRun
    metadata:
    generateName: dracon-github-com-kubernetes-kubernetes-
    namespace: default
    spec:
    serviceAccountName: dracon
    pipelineRef:
        name: dracon-github-com-kubernetes-kubernetes
    params:
    - name: repository_url
        value: https://github.com/kubernetes/kubernetes.git
    - name: consumer-elasticsearch-url
        value: http://quickstart-es-http:9200
    workspaces:
    - name: source-code-ws
        subPath: source-code
        volumeClaimTemplate:
        spec:
            accessModes:
            - ReadWriteOnce
            resources:
            requests:
                storage: 1Gi
    ```

4. Create the PipelineRun resource:

    ```bash
    $ kubectl apply -f pipelinerun.yaml
    ```

5. Observe the PipelineRun at http://localhost:8001/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/about
