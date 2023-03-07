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
    # Use `kubectl port-forward ...` so you can access Kubernetes services on your local machine.
    $ kubectl -n tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
    # Tekton Dashboard is now available at: http://localhost:9097
    ```

4. Create a Postgres DB to use as the Dracon deduplication DB:

    ```bash
    # Create a StatefulSet and Service for the Dracon deduplication DB. In production, we recommend using a production-ready or managed Postgres deployment.
    $ kubectl apply -f https://raw.githubusercontent.com/ocurity/dracon-community-pipelines/main/resources/deduplication-enricher-db.yaml
    ```

5. Install ECK and create an Elasticsearch + Kibana Dashboards. For more info, see [official documentation](https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-deploy-eck.html).

    ```bash
    # Create ECK CRDs.
    $ kubectl create -f https://download.elastic.co/downloads/eck/2.6.1/crds.yaml
    # Apply ECK operator resources.
    $ kubectl apply -f https://download.elastic.co/downloads/eck/2.6.1/operator.yaml
    # Create Elasticsearch.
    $ kubectl apply -f https://raw.githubusercontent.com/ocurity/dracon-community-pipelines/main/resources/eck-elasticsearch.yaml
    # Create Kibana.
    $ kubectl apply -f https://raw.githubusercontent.com/ocurity/dracon-community-pipelines/main/resources/eck-kibana.yaml
    # Use `kubectl port-forward ...` to access the Kibana UI:
    $ kubectl port-forward svc/quickstart-kb-http 5601:5601
    # You can obtain the password by examining the `quickstart-es-elastic-user` secret:
    # The username is `elastic`.
    $ kubectl get secret quickstart-es-elastic-user \
        -o=jsonpath='{.data.elastic}' \
        | base64 -d - \
        | xargs echo "$1"
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
    # Run `kubectl create ...` with this file.
    apiVersion: tekton.dev/v1beta1
    kind: PipelineRun
    metadata:
      generateName: dracon-github-com-kubernetes-kubernetes-
    spec:
      pipelineRef:
        name: dracon-github-com-kubernetes-kubernetes
      params:
        - name: repository_url
          value: https://github.com/kubernetes/kubernetes.git
        - name: consumer-elasticsearch-url
          value: http://quickstart-es-http.default.svc:9200
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
    $ kubectl create -f pipelinerun.yaml
    ```

5. Observe the PipelineRun at http://localhost:8001/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/about

6. Once the PipelineRun has finished, you can view the output in Kibana at http://localhost:5601.
