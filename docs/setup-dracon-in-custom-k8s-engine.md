# Setup Dracon on another Kubernetes engine (Not recommended)

If you don't want to use [KiND](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
locally, you can follow these steps to deploy Dracon and it's dependencies.

If you follow the [Getting Started](./getting-started.md) guide,
all these steps are taken care of by the `make install` command.

## Deploying Dracon dependencies

Dracon dependencies are split into two categories:

1. Dependencies that the system can't work without
2. Dependencies that the system doesn't need but are probably needed by most
   pipelines.

The dependency that Dracon can't function without is [Tekton](https://tekton.dev/)
and for many users
it is a good idea to deploy the Tekton Dashboard too for better visibility into
what's happening on the cluster. We offer a simple way of deploying these along
with an Nginx ingress controller with the command:

```bash
make dev-infra
```

## Deploying Dracon Pipeline dependencies using Helm packages

1. Deploy the Helm packages

> :warning: **Warning 2:** make sure that you have all the needed tools
> listed in the previous section installed in your system

For Dracon pipelines to run, they usually require the following services:

1. MongoDB
2. Elasticsearch
3. Kibana
4. MongoDB
5. Postgres

We use the Elastic Operator to spin up managed instances of Elasticsearch and
Kibana and the bitnami charts to deploy instances of PostgreSQL and MongoDB.

If you run the command:

```bash
make dev-dracon
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
  base64 -d && \
  echo
```

3. Expose the Kibana Dashboard

```bash
 # Use `kubectl port-forward ...` to access the Kibana UI:
 kubectl -n dracon port-forward svc/dracon-kb-kibana-kb-http 5601:5601
```

The username/password is the same as Kibana

## Deploy Dracon components

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
  --values deploy/dracon/values/dev.yaml \
  dracon-oss-components \
  oci://ghcr.io/ocurity/dracon/charts/dracon-oss-components 
```

## Applying manually migrations

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
```
