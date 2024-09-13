# Developer vars
# The following variables are used to define the developer environment
# e.g. what are the test packages, or the latest tag, these are used by make targets that build things
component_binaries=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/bin')
component_containers=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/docker')
component_containers_publish=$(component_containers:docker=publish)
protos=$(shell find . -not -path './vendor/*' -name '*.proto')
go_protos=$(protos:.proto=.pb.go)
latest_tag=$(shell git tag --list --sort="-version:refname" | head -n 1)
commits_since_latest_tag=$(shell git log --oneline $(latest_tag)..HEAD | wc -l)
GO_TEST_PACKAGES=$(shell go list ./... | grep -v /vendor/)

# Deployment vars
# The following variables are used to define the deployment environment
# e.g. what are the versions of the components, or the container registry, these are used by make targets that deploy things
CONTAINER_REPO=ghcr.io/ocurity/dracon
SOURCE_CODE_REPO=https://github.com/ocurity/dracon
DRACON_DEV_VERSION=$(shell echo $(latest_tag)$$([ $(commits_since_latest_tag) -eq 0 ] || echo "-$$(git log -n 1 --pretty='format:%h')" )$$([ -z "$$(git status --porcelain=v1 2>/dev/null)" ] || echo "-dirty" ))
DRACON_VERSION=$(shell (echo $(CONTAINER_REPO) | grep -q '^ghcr' && echo $(latest_tag)) || echo $(DRACON_DEV_VERSION) )
DRACON_OSS_COMPONENTS_NAME=dracon-oss-components
DRACON_OSS_COMPONENTS_PACKAGE_URL=oci://ghcr.io/ocurity/dracon/charts/$(DRACON_OSS_COMPONENTS_NAME)

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

DOCKER=docker
PROTOC=protoc

export

########################################
############# BUILD TARGETS ############
########################################
.PHONY: components component-binaries cmd/draconctl/bin protos build publish-component-containers publish-containers draconctl-image draconctl-image-publish clean-protos clean

$(component_binaries):
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 ./scripts/build_component_binary.sh $@

component-binaries: $(component_binaries)

$(component_containers): %/docker: %/bin
	./scripts/build_component_container.sh $@

components: $(component_containers)

cmd/draconctl/bin:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cmd/draconctl cmd/draconctl/main.go

draconctl-image: cmd/draconctl/bin
	$(DOCKER) build --platform linux/arm64 -t "${CONTAINER_REPO}/draconctl:${DRACON_VERSION}" \
		$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
		-f containers/Dockerfile.draconctl .

draconctl-image-publish: draconctl-image
	$(DOCKER) push "${CONTAINER_REPO}/draconctl:${DRACON_VERSION}"

third_party/tektoncd/swagger-v$(TEKTON_VERSION).json:
	@wget "https://raw.githubusercontent.com/tektoncd/pipeline/v$(TEKTON_VERSION)/pkg/apis/pipeline/v1beta1/swagger.json" -O $@

components/base/openapi_schema.json: third_party/tektoncd/swagger-v$(TEKTON_VERSION).json
	@cp $< $@

$(go_protos): %.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative $<

protos: $(go_protos)

build: components protos
	@echo "done building"

$(component_containers_publish): %/publish: %/docker
	./scripts/publish_component_container.sh $@

publish-component-containers: $(component_containers_publish)

publish-containers: publish-component-containers draconctl-image-publish

clean-protos:
	@find . -not -path './vendor/*' -name '*.pb.go' -delete

clean-migrations-compose:
	cd tests/migrations/ && docker compose rm --force

clean: clean-protos clean-migrations-compose

########################################
######### CODE QUALITY TARGETS #########
########################################
.PHONY: lint install-lint-tools tests go-tests fmt fmt-proto fmt-go install-go-fmt-tools migration-tests

lint:
# we need to redirect stderr to stdout because Github actions don't capture the stderr lolz
	@reviewdog -fail-on-error -diff="git diff origin/main" -filter-mode=added 2>&1

install-lint-tools:
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/mgechev/revive@latest
	@go install github.com/sivchari/containedctx/cmd/containedctx@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/polyfloyd/go-errorlint@latest
	@go install github.com/kisielk/errcheck@latest
	@go install github.com/rhysd/actionlint/cmd/actionlint@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@go install github.com/bufbuild/buf/cmd/buf@v1.32.2
	@npm ci

install-go-test-tools:
	@go install gotest.tools/gotestsum@latest

go-tests:
	@mkdir -p tests/output
	@gotestsum --junitfile tests/output/unit-tests.xml -- -race -coverprofile tests/output/cover.out $(GO_TEST_PACKAGES)

go-cover: go-tests
	@go tool cover -html=tests/output/cover.out -o=tests/output/cover.html && open tests/output/cover.html

migration-tests: cmd/draconctl/bin
	cd tests/migrations/ && docker compose up --abort-on-container-exit --build --exit-code-from tester

test: go-tests migration-tests

fmt-proto:
	@echo "Tidying up Proto files"
	@buf format -w ./api/proto

install-go-fmt-tools:
	@go install github.com/bufbuild/buf/cmd/buf@v1.28.1
	@go install golang.org/x/tools/cmd/goimports@latest

fmt-go:
	@echo "Tidying up Go files"
	@gofmt -l -w $$(find . -name *.go -not -path "./vendor/*" | xargs -n 1 dirname | uniq)
	@goimports -local $$(cat go.mod | grep -E "^module" | sed 's/module //') -w $$(find . -name *.go -not -path "./vendor/*" | xargs -n 1 dirname | uniq)

install-md-fmt-tools:
	@npm ci

fmt-md:
	@echo "Tidying up MD files"
	@npm run format

fmt: fmt-go fmt-proto fmt-md

########################################
########## DEBUGGING TARGETS ###########
########################################

print-%:
	@echo $($*)

########################################
########## DEPLOYMENT TARGETS ##########
########################################
.PHONY: deploy-nginx deploy-arangodb-crds deploy-arangodb-operator add-es-helm-repo deploy-elasticoperator \
		tektoncd-dashboard-helm deploy-tektoncd-dashboard add-bitnami-repo dev-dracon dev-deploy dev-teardown \
		install install-oss-components deploy-cluster

deploy-nginx:
	@helm upgrade nginx-ingress https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-$(NGINX_INGRESS_VERSION)/ingress-nginx-$(NGINX_INGRESS_VERSION).tgz \
		--install \
		--namespace $(NGINX_INGRESS_NS) \
		--create-namespace \
		--set "controller.admissionWebhooks.enabled=false"

deploy-arangodb-crds:
	@helm upgrade arangodb-crds https://github.com/arangodb/kube-arangodb/releases/download/$(ARANGODB_VERSION)/kube-arangodb-crd-$(ARANGODB_VERSION).tgz \
		--install

deploy-arangodb-operator:
	@helm install --generate-name https://github.com/arangodb/kube-arangodb/releases/download/1.2.40/kube-arangodb-1.2.40.tgz

add-es-helm-repo:
	@helm repo add elastic https://helm.elastic.co
	@helm repo update

deploy-elasticoperator: add-es-helm-repo
	@helm upgrade elastic-operator elastic/eck-operator \
		--install \
		--namespace $(ES_NAMESPACE) \
		--create-namespace \
		--version=$(ES_OPERATOR_VERSION)

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

add-bitnami-repo:
	@helm repo add bitnami https://charts.bitnami.com/bitnami

deploy-cluster:
	@scripts/kind-with-registry.sh

install: deploy-cluster dev-infra deploy-elasticoperator deploy-arangodb-crds add-bitnami-repo
	@echo "fetching dependencies if needed"
	@helm dependency build ./deploy/dracon/chart

	@echo "deploying dracon"
	@helm upgrade dracon ./deploy/dracon/chart \
	 	  --install \
		  --values ./deploy/dracon/values/dev.yaml \
		  --create-namespace \
		  --set "image.registry=$(CONTAINER_REPO)" \
		  --namespace $(DRACON_NS) \
		  --version $(DRACON_VERSION) \
		  --wait

	@echo "Applying migrations"
	@helm upgrade deduplication-db-migrations ./deploy/deduplication-db-migrations/chart \
		  --install \
		  --values ./deploy/deduplication-db-migrations/values/dev.yaml \
		  --create-namespace \
		  --set "image.registry=$(CONTAINER_REPO)" \
		  --namespace $(DRACON_NS) \
		  --set "image.tag=$(DRACON_VERSION)" \
		  --wait

	@echo "Installing Components"
	# we are setting the container repo to it's own value so that we can override it from other make targets
	# e.g. when installing oss components from locally built components, we want to `make install` with CONTAINER_REPO being the kind-registry, and the package_url being the component tar.gz
	$(MAKE) install-oss-components CONTAINER_REPO=$(CONTAINER_REPO) DRACON_OSS_COMPONENTS_PACKAGE_URL=$(DRACON_OSS_COMPONENTS_PACKAGE_URL)

dev-deploy-oss-components:
	@echo "Deploying components in local dracon instance"
	$(MAKE) dev-build-oss-components CONTAINER_REPO=$(CONTAINER_REPO)
	$(MAKE) install-oss-components CONTAINER_REPO=$(CONTAINER_REPO) DRACON_OSS_COMPONENTS_PACKAGE_URL=$(DRACON_OSS_COMPONENTS_PACKAGE_URL)

install-oss-components:
	@helm upgrade $(DRACON_OSS_COMPONENTS_NAME) \
		$(DRACON_OSS_COMPONENTS_PACKAGE_URL) \
		--install \
		--create-namespace \
		--namespace $(DRACON_NS) \
		--set image.registry=$(CONTAINER_REPO) \
		--values ./deploy/deduplication-db-migrations/values/dev.yaml
	@echo "Done! Bumped version to $(DRACON_VERSION)"

dev-build-oss-components: cmd/draconctl/bin
	@echo "Building open-source components for local dracon instance..."
	$(eval CONTAINER_REPO:=localhost:5000)

	$(MAKE) -j 16 publish-component-containers CONTAINER_REPO=$(CONTAINER_REPO)
	@./bin/cmd/draconctl components package \
		--version $(DRACON_VERSION) \
		--chart-version $(DRACON_VERSION) \
		--name $(DRACON_OSS_COMPONENTS_NAME) \
		./components

dev-dracon:
	$(eval CONTAINER_REPO:=localhost:5000)
	$(eval DRACON_OSS_COMPONENTS_PACKAGE_URL:=./$(DRACON_OSS_COMPONENTS_NAME)-$(DRACON_VERSION).tgz)
	$(eval IN_CLUSTER_CONTAINER_REPO:=kind-registry:5000)
	
	$(MAKE) -j 16 publish-containers CONTAINER_REPO=$(CONTAINER_REPO)
	$(MAKE) -j 16 dev-build-oss-components CONTAINER_REPO=$(CONTAINER_REPO)

	$(MAKE) install CONTAINER_REPO=$(IN_CLUSTER_CONTAINER_REPO) DRACON_OSS_COMPONENTS_PACKAGE_URL=$(DRACON_OSS_COMPONENTS_PACKAGE_URL)

dev-infra: deploy-nginx deploy-tektoncd-pipeline deploy-tektoncd-dashboard

dev-deploy: deploy-cluster dev-infra dev-dracon

dev-teardown:
	@kind delete clusters dracon-demo

generate-protos: install-lint-tools
	@echo "Generating Protos"
	@buf generate
