component_binariess=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/bin')
component_containers=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/docker')
component_kustomizations=$(shell find ./components -name kustomization.yaml | xargs -I'{}' sh -c 'echo $$(dirname {})/kustomization')
component_containers_publish=$(component_containers:docker=publish)
example_pipeline_kustomizations=$(shell find ./examples/pipelines -name kustomization.yaml | xargs -I'{}' sh -c 'echo $$(dirname {})/helm-templates')
protos=$(shell find . -not -path './vendor/*' -name '*.proto')
go_protos=$(protos:.proto=.pb.go)
latest_tag=$(shell git tag --list --sort="-version:refname" | head -n 1)
commits_since_latest_tag=$(shell git log --oneline $(latest_tag)..HEAD | wc -l)

GO_TEST_PACKAGES=$(shell go list ./... | grep -v /vendor/)
CONTAINER_REPO=europe-west1-docker.pkg.dev/oc-dracon-saas/demo/ocurity/dracon
DRACON_VERSION=$(shell echo $(latest_tag)$$([ $(commits_since_latest_tag) -eq 0 ] || echo "-$$(git log -n 1 --pretty='format:%h')" )$$([ -z "$$(git status --porcelain=v1 2>/dev/null)" ] || echo "-dirty" ))
TEKTON_VERSION=0.44.0
TEKTON_DASHBOARD_VERSION=0.29.2
ARANGODB_VERSION=1.2.19
NGINX_INGRESS_VERSION=4.2.5
NGINX_INGRESS_NS=ingress-nginx
NAMESPACE=default
ES_NAMESPACE=elastic-system
ES_OPERATOR_VERSION=2.2.0
ES_VERSION=8.3.2
MONGODB_VERSION=13.3.0
PG_VERSION=11.9.8
DRACON_NS=dracon
TEKTON_NS=tekton-pipelines
ARANGODB_NS=arangodb
BASE_IMAGE=scratch

DOCKER=docker
PROTOC=protoc

export

########################################
############# BUILD TARGETS ############
########################################
.PHONY: components component-binaries protos build publish-component-containers clean-protos clean

$(component_binariess):
	./scripts/build_component_binary.sh $@

component-binaries: $(component_binariess)

$(component_containers): %/docker: %/bin
	./scripts/build_component_container.sh $@

components: $(component_containers)

bin/cmd/draconctl:
	@go build -o bin/cmd/draconctl cmd/draconctl/main.go

third_party/tektoncd/swagger-v$(TEKTON_VERSION).json:
	@wget "https://raw.githubusercontent.com/tektoncd/pipeline/v$(TEKTON_VERSION)/pkg/apis/pipeline/v1beta1/swagger.json" -O $@

components/base/openapi_schema.json: third_party/tektoncd/swagger-v$(TEKTON_VERSION).json
	@cp $< $@

$(go_protos): %.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative $<

protos: $(go_protos)

build: components protos example-pipelines-helm-manifests
	@echo "done building"

$(component_containers_publish): %/publish: %/docker
	./scripts/publish_component_container.sh $@

publish-component-containers: $(component_containers_publish)

$(example_pipeline_kustomizations): bin/cmd/draconctl
	@echo "Generating templates for $$(dirname $@)"
	@mkdir -p "$$(dirname $@)/templates"
	@bin/cmd/draconctl pipelines build --out "$$(dirname $@)/templates/all.yaml" $$(dirname $@)

example-pipelines-helm-manifests: $(example_pipeline_kustomizations)

clean-protos:
	@find . -not -path './vendor/*' -name '*.pb.go' -delete

clean-migrations-compose:
	cd tests/migrations/ && docker compose rm --force

clean: clean-protos clean-migrations-compose

########################################
######### CODE QUALITY TARGETS #########
########################################
.PHONY: lint install-lint-tools tests go-tests fmt fmt-proto fmt-go migration-tests

lint:
	@reviewdog -fail-on-error $$([ "${CI}" = "true" ] && echo "-reporter=github-pr-review") -diff="git diff origin/main" -tee

install-lint-tools:
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/mgechev/revive@latest
	@go install github.com/sivchari/containedctx/cmd/containedctx@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/polyfloyd/go-errorlint@latest
	@go install github.com/kisielk/errcheck@latest
	@go install github.com/rhysd/actionlint/cmd/actionlint@latest
	@go install github.com/bufbuild/buf/cmd/buf@v1.28.1

go-tests:
	@go test -race -json $(GO_TEST_PACKAGES)

migration-tests: bin/cmd/draconctl
	cd tests/migrations/ && docker compose up --abort-on-container-exit --build --exit-code-from tester

test: go-tests migration-tests

fmt-proto:
	@echo "Tidying up Proto files"
	@buf format -w ./api/proto

fmt-go:
	@echo "Tidying up Go files"
	@gofmt -l -w $$(find . -name *.go -not -path "./vendor/*" | xargs -n 1 dirname | uniq)

fmt: fmt-go fmt-proto

########################################
########## DEPLOYMENT TARGETS ##########
########################################
.PHONY: deploy-arangodb-crds deploy-arangodb dev-deploy deploy-elasticsearch deploy-mongodb deploy-pg deploy-tektoncd-pipeline tektoncd-pipeline-helm tektoncd-dashboard-helm

deploy-arangodb-crds:
	@helm upgrade arangodb-crds https://github.com/arangodb/kube-arangodb/releases/download/$(ARANGODB_VERSION)/kube-arangodb-crd-$(ARANGODB_VERSION).tgz \
		--install

deploy-arangodb: deploy-arangodb-crds
	@helm upgrade arangodb-instance deploy/arangodb/ \
		--install \
		--namespace $(ARANGODB_NS) \
		--create-namespace \
		--values=deploy/arangodb/values.yaml

deploy-nginx:
	@helm upgrade nginx-ingress https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-$(NGINX_INGRESS_VERSION)/ingress-nginx-$(NGINX_INGRESS_VERSION).tgz \
		--install \
		--namespace $(NGINX_INGRESS_NS) \
		--create-namespace \
		--set "controller.admissionWebhooks.enabled=false"

add-es-helm-repo:
	@helm repo add elastic https://helm.elastic.co
	@helm repo update

deploy-elasticoperator: add-es-helm-repo
	@helm upgrade elastic-operator elastic/eck-operator \
		--install \
		--namespace $(ES_NAMESPACE) \
		--create-namespace \
		--version=$(ES_OPERATOR_VERSION)

deploy-elasticsearch: deploy-elasticoperator
	@helm upgrade dracon-es deploy/elasticsearch/ \
		--install \
		--set version=$(ES_VERSION) \
		--namespace $(DRACON_NS) \
		--create-namespace

deploy-kibana: deploy-elasticsearch
	@helm upgrade dracon-kb deploy/kibana/ \
		--install \
		--set version=$(ES_VERSION) \
		--set es_name=dracon-es-elasticsearch \
		--namespace $(DRACON_NS) \
		--version $(ES_VERSION)

deploy-mongodb:
	@helm upgrade consumer-mongodb https://charts.bitnami.com/bitnami/mongodb-$(MONGODB_VERSION).tgz \
		--install \
		--namespace $(DRACON_NS) \
		--create-namespace

deploy-pg:
	@helm upgrade pg https://charts.bitnami.com/bitnami/postgresql-$(PG_VERSION).tgz \
		--install \
		--namespace $(DRACON_NS) \
		--create-namespace \
		--values=deploy/enrichment-db/values.yaml

deploy/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml:
	@wget "https://storage.googleapis.com/tekton-releases/pipeline/previous/v$(TEKTON_VERSION)/release.yaml" -O $@

tektoncd-pipeline-helm: deploy/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml
	./scripts/generate_tektoncd_pipeline_helm.sh $(TEKTON_VERSION)

deploy-tektoncd-pipeline: tektoncd-pipeline-helm
	@helm upgrade tektoncd ./deploy/tektoncd/pipeline \
		--install \
		--namespace $(TEKTON_NS) \
		--create-namespace

deploy/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml:
    @wget "https://github.com/tektoncd/dashboard/releases/download/v$(TEKTON_DASHBOARD_VERSION)/tekton-dashboard-release.yaml" -O $@

tektoncd-dashboard-helm: deploy/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml
	./scripts/generate_tektoncd_dashboard_helm.sh $(TEKTON_DASHBOARD_VERSION)

deploy-tektoncd-dashboard: tektoncd-dashboard-helm
	@helm upgrade tektoncd-dashboard ./deploy/tektoncd/dashboard \
		--install \
		--values ./deploy/tektoncd/dashboard/values.yaml \
		--namespace $(TEKTON_NS)

dev-deploy: deploy-nginx deploy-arangodb deploy-kibana deploy-mongodb deploy-pg deploy-tektoncd-pipeline deploy-tektoncd-dashboard
