# Dracon Demos

This tutorial will show you how to run the
[Demo Dracon Pipelines](examples/pipelines/golang-project) and see the results in Kibana at the end.
The first section describes how to setup locally a KinD cluster named `dracon-demo` to quickly demo
Tekton. Then we describe how to setup the services and finally how to deploy the example code.

## Prerequisites

- KinD [official documentation](https://kind.sigs.k8s.io/docs/user/quick-start/#creating-a-cluster):
- Helm []

## KinD cluster setup

1. Create a local Docker registry:

```bash
$ docker inspect -f '{{.State.Running}}' "kind-dracon-demo-registry" 2>/dev/null || \
   docker run -d --restart=always -p "127.0.0.1:5001:5000" --network bridge --name "kind-dracon-demo-registry" registry:2
```

2. Create KinD cluster
```bash
$ cat <<EOF | kind create cluster --name dracon-demo --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "/etc/containerd/certs.d"
EOF
```

3. Connect KinD cluster to local Docker registry
```bash
$ REGISTRY_DIR="/etc/containerd/certs.d/localhost:5001"
$ for node in $(kind get nodes); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://kind-dracon-demo-registry:5000"]
EOF
done
$ [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' kind-dracon-demo-registry) = 'null')" ] || docker network connect kind kind-dracon-demo-registry
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:5001"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
```

## Building components and deploying Tekton

1. Building all the components is fairly simple. Checkout the version of Tekton you want to use
and run the following:

```bash
$ make components
```

* If you used KinD with a local repository run the following:

```bash
$ make components DOCKER_REPO=localhost:5001/ocurity/dracon
```

2. Deploy Tekton

```bash
$ make dev-deploy
```

> :warning: **Warning 1:** make sure that all pods are up an running before proceeding

3. Deploy the Helm chart with the golang project

```bash
$ helm upgrade golang-project-pipeline ./examples/pipelines/golang-project/ \
    --install \
    --namespace dracon
```

4. Start a pipeline run

```bash
$ kubectl apply -n dracon -f ./examples/pipelines/golang-project/pipelinerun/pipelinerun.yaml
```

5. Expose the Tekton dashboard and Kibana UI

```bash
$ kubectl -n tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
$ kubectl -n dracon port-forward svc/dracon-kb-kibana-kb-http 5601:5601
```

6. Get the token to access Kibana

```bash
# The username is `elastic`.
$ kubectl -n dracon get secret dracon-es-elasticsearch-es-elastic-user \
        -o=jsonpath='{.data.elastic}' \
        | base64 -d - \
        | xargs echo "$1"
```

7. Deploy the Golang pipeline



8. Wait for the pipeline to finish running by monitoring it in http://localhost:9097.

9.  Once the pipelinerun has finished running you can view your results in Kibana: http://localhost:5601.
