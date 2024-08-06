# Getting Started with Dracon

This guide will help to quickly setup Dracon on a Kubernetes cluster and get a
pipeline running. The first step is to create a dev Kubernetes cluster in order
to deploy Tekton. We suggest you use KiND to provision a local test cluster
quickly. If you already have a K8s cluster then you can skip directly to the
[Deploying Dracon dependencies](#deploying-dracon-dependencies) section.

We support two ways of deploying Dracon.

1. Using the Helm packages that we distribute
2. Using your local copy of this repository

The 2nd option is useful when you are developing components or new functionality

## Tools you will need

You will need to have the following tools installed in your system:

1. [KiND](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
2. [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)
3. [Docker](https://docs.docker.com/engine/install/)
4. [Helm](https://helm.sh/docs/intro/install/)

## Setting up a [KinD](https://kind.sigs.k8s.io/) cluster

KinD is is a tool for running local Kubernetes clusters using Docker container
“nodes”.

Create a KinD cluster named `dracon-demo` with its own Docker registry. You can
use our Bash script or you can check for more info the
[official documentation](https://kind.sigs.k8s.io/docs/user/quick-start/#creating-a-cluster):

```bash
./scripts/kind-with-registry.sh
```

> :warning: **Warning 1:** make sure you can connect to the KiND control plane
> before proceeding. some times it takes a bit for everything to get started.

## Deploying Dracon dependencies

Dracon dependencies are split into two categories:

1. Dependencies that the system can't work without
2. Dependencies that the system doesn't need but are probably needed by most
   pipelines.

The dependency that Dracon can't function without is Tekton and for many users
it is a good idea to deploy the Tekton Dashboard too for better visibility into
what's happening on the cluster. We offer a simple way of deploying these along
with an Nginx ingress controller with the command:

```bash
make dev-infra
```

## Deploying Dracon

### Deploying Dracon Pipeline dependencies using Helm packages

1. Deploy the Helm packages

> :warning: **Warning 2:** make sure that you have all the needed tools
> listed in the previous section installed in your system

For Dracon pipelines to run, they usually require the following services:

1. MongoDB
2. Elastic Search
3. Kibana
4. MongoDB
5. Postgres

We use the Elastic Operator to spin up managed instances of Elasticsearch and
Kibana and the bitnami charts to deploy instances of PostgreSQL and MongoDB.

If you run the command:

```bash
make dev-dracon DRACON_VERSION=v0.13.0
```

You will deploy the Elastic Operator on the cluster and the Dracon Helm
package. Depending on the capabilities of your workstation this will probably
take a couple of minutes, it's perfect time to go get a cup of coffee ;).

```text
   )  (
  (   ) )
   ) ( (
  -------
.-\     /
'- \   /
  _______
```

`espresso cup by @ptzianos`

The Dracon Helm package lists as dependencies the Bitnami charts for Postgres
and MongoDB. The values used are in the `deploy/dracon/values/dev.yaml` file.

1. Expose the TektonCD Dashboard

   ```bash
     kubectl -n tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
   ```

2. Expose the Kibana Dashboard.

   ```bash
   # Use `kubectl port-forward ...` to access the Kibana UI:
   kubectl -n dracon port-forward svc/dracon-kb-kibana-kb-http 5601:5601
   # You can obtain the password by examining the 
   # `dracon-es-elasticsearch-es-elastic-user` secret:
   # The username is `elastic`.
   kubectl -n dracon get secret dracon-es-elasticsearch-es-elastic-user \
            -o=jsonpath='{.data.elastic}' | \
            base64 -d &&\
            echo
   ```

3. Expose the Kibana Dashboard

   ```bash
   # Use `kubectl port-forward ...` to access the Kibana UI:
   kubectl -n dracon port-forward svc/dracon-kb-kibana-kb-http 5601:5601
   ```

   The username/password is the same as Kibana

### Deploy Dracon components

The components that are used to build our pipelines are comprised out of two
pieces:

1. a wrapper around the binary of the tool that we wish to execute packaged
   into a container.
2. a Tekton Task file that describes how to execute the component.

We provide Helm packages with all our components that can be easily installed
as follows:

```bash
helm upgrade \
  --install \
  --namespace dracon \
  --version 0.8.0
  dracon-oss-components\
  oci://ghcr.io/ocurity/dracon/charts/dracon-oss-components 
```

### Deploying a custom version of Dracon components

So, the first step is to build all the containers and push them to a registry
that your cluster has access to. We use `make` to package our containers. For
each component our Make will automatically generate a phony target with the
path `components/{component type}/{component name}/docker`. We have a top-level
target that creates all the component containers along with a couple of extra
containers our system uses, such as draconctl.

The following examples are using the local container registry used by the KiND
cluster, but make sure that you replace the URL with the registry URL that you
are using, if you are using something else:

```bash
make publish-component-containers CONTAINER_REPO=localhost:5000/ocurity/dracon
```

\* Notice that the repo we are using is slightly different than the
one we pushed the images in the previous step. That's because with local
registries the registry is exposed on a port in localhost, however inside the
KiND cluster, that's not the case. Instead the registry's host is
`kind-registry:5000`. This is also going to be important later when we will
deploy the pipelines and their image repositories will also have to be set to
this value.

\*\*Make sure that you use the draconctl image that you pushed in the repository

#### Using a different base image for your images

If you need your images to have a different base image then you can pass the
`BASE_IMAGE` variable to the `components` or `publish-component-containers` to
change it to whatever you need. The targets build the binaries and place them in
the `bin` directory and then other targets package them into containers with
`scratch` as the base image.

There are some components that require extra components or special treatment and
these components have their own Makefiles. In those cases you can place a
`.custom_image` file in the directory with the base image you wish to use and
that will be picked up by the Makefile and build the container.

#### Deploying your custom Dracon components Helm package

You can package your components into a Helm package by running the following
command:

```bash
export CUSTOM_DRACON_VERSION=$(make print-DRACON_VERSION)
export CUSTOM_HELM_COMPONENT_PACKAGE_NAME=
make cmd/draconctl/bin
bin/cmd/draconctl components package \
  --version ${CUSTOM_DRACON_VERSION} \
  --chart-version ${CUSTOM_DRACON_VERSION} \
  --name ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME} \
  ./components
helm upgrade ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME} ./${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}-${CUSTOM_DRACON_VERSION}.tgz \
  --install \
  --namespace dracon
```

If your custom components are local, you need to override the component registry
you can do so with the following slightly modified helm command

```bash
helm upgrade ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME} ./${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}-${CUSTOM_DRACON_VERSION}.tgz \
  --install \
  --namespace dracon\
  --set image.registry=kind-registry:5000/ocurity/dracon
```

After changes to your components you need to redeploy, you can do so as such:

```bash

export CUSTOM_DRACON_VERSION=$(make print-DRACON_VERSION)
make publish-component-containers CONTAINER_REPO=localhost:5000/ocurity/dracon
bin/cmd/draconctl components package   --version ${CUSTOM_DRACON_VERSION}   \
  --chart-version ${CUSTOM_DRACON_VERSION}   \
  --name ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}   \
  ./components 
helm upgrade ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME} \
  ./${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}-${CUSTOM_DRACON_VERSION}.tgz   \
  --install \
  --namespace dracon \
  --set image.registry=kind-registry:5000/ocurity/dracon
```

````
### Applying manually migrations

There some migrations that should be applied to the postgres instance so that
the enrichment components can store and retrieve data from it. In order to apply
the migrations you need to run the following command (the container with the
`draconctl` binary and the migration scripts was built and pushed in the
previous step):

```bash
kubectl apply -n dracon -f deploy/dracon/serviceaccount.yaml
kubectl apply -n dracon -f deploy/dracon/role.yaml
kubectl apply -n dracon -f deploy/dracon/rolebinding.yaml
make cmd/draconctl/bin

export DRACONCTL_MIGRATIONS_PATH='/etc/dracon/migrations/enrichment'
bin/cmd/draconctl migrations apply \
  --namespace dracon \
  --as-k8s-job \
  --image "${CONTAINER_REPO}/draconctl:${CUSTOM_DRACON_VERSION}" \
  --url "postgresql://dracon:dracon@dracon-enrichment-db.dracon.svc.cluster.local?sslmode=disable" \
````

## Running one of the example pipelines

You can easily check Dracon in action if you run one of the example pipelines
included in the repo.
Running the `golang-project` is as simple as running:

```bash
make cmd/draconctl/bin
bin/cmd/draconctl pipelines deploy ./examples/pipelines/golang-project
```

Then you can run an instance of the pipeline as follows:

```bash
kubectl create \
  -n dracon \
  -f ./examples/pipelines/golang-project/pipelinerun.yaml
```
