.PHONY: build clean clean-protos lint fmt fmt-go fmt-proto tests go-tests release_notes kustomizations

proto_defs=$(shell find . -name "*.proto" -not -path "./vendor/*")
go_protos=$(proto_defs:.proto=.pb.go)
component_binariess=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/bin')
component_containers=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/docker')
component_containers_push=$(component_containers:docker=push)
component_containers_retag=$(component_containers:docker=retag)
component_tasks=$(shell find ./components -name task.yaml -not -path "./components/sources/git/*" -not -path "./components/enrichers/aggregator/*" -not -path "./components/producers/aggregator/*" -not -path "./components/base/*")
component_kustomizations=$(component_tasks:task.yaml=kustomization.yaml)

GO_TEST_PACKAGES=./...

DOCKER_REPO=ghcr.io/ocurity/dracon
TEKTON_VERSION=0.44.0
TEKTON_DASHBOARD_VERSION=0.29.2
ARANGODB_VERSION=1.2.19
NGINX_INGRESS_VERSION=4.2.5
NAMESPACE=default
ES_NAMESPACE=elastic-system
ES_OPERATOR_VERSION=2.2.0
MONGODB_VERSION=13.3.0
PG_VERSION=11.9.8
NGINX_INGRESS_NS=ingress-nginx
DRACON_NS=dracon
TEKTON_NS=tekton-pipelines
ARANGODB_NS=arangodb

latest_tag=$(shell git tag --list --sort="-version:refname" | head -n 1)
commits_since_latest_tag=$(shell git log --oneline $(latest_tag)..HEAD | wc -l)
DRACON_VERSION=$(shell echo $(latest_tag)$$([ $(commits_since_latest_tag) -eq 0 ] || echo "-$$(git log -n 1 --pretty='format:%h')" )$$([ -z "$$(git status --porcelain=v1 2>/dev/null)" ] || echo "-dirty" ))

GO=go
PROTOC=protoc
DOCKER=docker

build: $(go_protos)
	@echo "done building"

$(go_protos): %.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative $<

clean-kustomizations:
	rm -rf $(component_kustomizations)

clean-protos:
	rm -rf $(go_protos)

clean: clean-protos clean-kustomizations

fmt-proto:
	@echo "Tidying up Proto files"
	@buf format -w "./api/proto"

fmt-go:
	@echo "Tidying up Go files"
	@gofmt -l -w $$(find . -name *.go -not -path "./vendor/*" | xargs -n 1 dirname | uniq)

fmt: fmt-go fmt-proto

lint:
	@if [ "${CI}" = "true" ]; then\
		reviewdog -fail-on-error -reporter=github-pr-review;\
	else\
		reviewdog -fail-on-error -diff="git diff origin/main";\
	fi
	@golangci-lint

go-tests:
	go test -race -json $(GO_TEST_PACKAGES)

tests: go-tests

$(component_binariess):
	$(shell                                                                                         \
        set -e;                                                                                     \
        EXECUTABLE=$$(basename $$(dirname $@));                                                     \
        echo "$@" | grep -Eq ^components/producers/.*$ && EXECUTABLE=$${EXECUTABLE}-parser || true; \
        EXECUTABLE_FOLDER=$$(dirname $@);                                                           \
        EXECUTABLE_PATH=$$(dirname $$(dirname $@))/$${EXECUTABLE};                                  \
		echo "building bin/$${EXECUTABLE_PATH}" > /dev/stderr;                                      \
        $(GO) build -o bin/$${EXECUTABLE_PATH} ./$${EXECUTABLE_FOLDER}/main.go;                     \
	)

component-binaries: $(component_binariess)

$(component_containers): %/docker: %/bin
	$(shell                                                                                                                                                                      \
        set -e;                                                                                                                                                                  \
        EXECUTABLE=$$(basename $$(dirname $@));                                                                                                                                  \
        echo "$@" | grep -Eq ^components/producers/.*$ && EXECUTABLE=$${EXECUTABLE}-parser || true;                                                                              \
        EXECUTABLE_FOLDER=$$(dirname $@);                                                                                                                                        \
        EXECUTABLE_PATH=$$(dirname $$(dirname $@))/$${EXECUTABLE};                                                                                                               \
        if $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --dry-run package >/dev/null 2>&1;                                                                              \
        then                                                                                                                                                                     \
            $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --quiet build DOCKER_REPO=$(DOCKER_REPO) DRACON_VERSION=$(DRACON_VERSION); EXECUTABLE_PATH=$(EXECUTABLE_PATH); \
        else                                                                                                                                                                     \
            DOCKERFILE_TEMPLATE="                                                                                                                                                \
                FROM scratch                                                                                                                                                  \n \
                COPY $${EXECUTABLE_PATH} /bin/$${EXECUTABLE}                                                                                                                  \n \
                ENTRYPOINT ["/bin/$${EXECUTABLE}"]                                                                                                                            \n \
            ";                                                                                                                                                                   \
            DOCKERFILE=$$(mktemp);                                                                                                                                               \
            printf "$${DOCKERFILE_TEMPLATE}" > $${DOCKERFILE};                                                                                                                   \
            docker build -t $(DOCKER_REPO)/$${EXECUTABLE_FOLDER}:$(DRACON_VERSION) -f $${DOCKERFILE} ./bin;                                                                      \
        fi;                                                                                                                                                                      \
        if $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --dry-run package-extras >/dev/null 2>&1;                                                                       \
        then                                                                                                                                                                     \
			$(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --quiet package-extras DOCKER_REPO=$(DOCKER_REPO) DRACON_VERSION=$(DRACON_VERSION);  fi; \
	)

components: $(component_containers)

$(component_containers_push):
	docker push $(DOCKER_REPO)/$$(dirname $@):$(DRACON_VERSION)

push-component-containers: $(component_containers_push)

$(component_containers_retag):
	docker tag $(OLD_DOCKER_REPO)/$$(dirname $@):$(DRACON_VERSION) $(DOCKER_REPO)/$$(dirname $@):$(DRACON_VERSION)

retag-component-containers: $(component_containers_retag)

third_party/tektoncd/swagger-v$(TEKTON_VERSION).json:
	wget "https://raw.githubusercontent.com/tektoncd/pipeline/v$(TEKTON_VERSION)/pkg/apis/pipeline/v1beta1/swagger.json" -O $@

api/openapi/tekton/openapi_schema.json: third_party/tektoncd/swagger-v$(TEKTON_VERSION).json
	./scripts/generate_openapi_schema.sh $< $@

components/base/openapi_schema.json: third_party/tektoncd/swagger-v$(TEKTON_VERSION).json
	cp $< $@

release_notes:
	git log --date=short --pretty='format:- %cd %s' -n 20

bin/kustomize-component-generator:
	go build -o bin/kustomize-component-generator build/tools/kustomize-component-generator/main.go

$(component_kustomizations): bin/kustomize-component-generator
	bin/kustomize-component-generator -task "$$(dirname $@)/task.yaml"

kustomizations: $(component_kustomizations)

print-%:
	@echo $($*)

.PHONY: deploy-arangodb-crds deploy-arangodb dev-deploy deploy-elasticsearch deploy-mongodb deploy-pg deploy-tektoncd-pipeline tektoncd-pipeline-helm tektoncd-dashboard-helm
deploy-arangodb-crds:
	helm upgrade arangodb-crds https://github.com/arangodb/kube-arangodb/releases/download/$(ARANGODB_VERSION)/kube-arangodb-crd-$(ARANGODB_VERSION).tgz \
		--install

deploy-arangodb: deploy-arangodb-crds
	helm upgrade arangodb-instance deploy/arangodb/ \
		--install \
		--namespace $(ARANGODB_NS) \
		--create-namespace \
		--values=deploy/arangodb/values.yaml

deploy-nginx:
	helm upgrade nginx-ingress https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-$(NGINX_INGRESS_VERSION)/ingress-nginx-$(NGINX_INGRESS_VERSION).tgz \
		--install \
		--namespace $(NGINX_INGRESS_NS) \
		--create-namespace \

add-es-helm-repo:
	helm repo add elastic https://helm.elastic.co
	helm repo update

deploy-elasticoperator: add-es-helm-repo
	helm upgrade elastic-operator elastic/eck-operator \
		--install \
		--namespace $(ES_NAMESPACE) \
		--create-namespace \
		--version=$(ES_OPERATOR_VERSION)

deploy-elasticsearch: deploy-elasticoperator
	helm upgrade elasticsearch deploy/elasticsearch/ \
		--install \
		--namespace $(DRACON_NS) \
		--create-namespace

deploy-mongodb:
	helm upgrade consumer-mongodb https://charts.bitnami.com/bitnami/mongodb-$(MONGODB_VERSION).tgz \
		--install \
		--namespace $(DRACON_NS) \
		--create-namespace

deploy-pg:
	helm upgrade pg https://charts.bitnami.com/bitnami/postgresql-$(PG_VERSION).tgz \
		--install \
		--namespace $(DRACON_NS) \
		--create-namespace \
		--values=deploy/enrichment-db/values.yaml
	
deploy/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml:
	wget "https://storage.googleapis.com/tekton-releases/pipeline/previous/v$(TEKTON_VERSION)/release.yaml" -O $@

tektoncd-pipeline-helm: deploy/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml
	$(shell                                                                                                                                                                                                       \
	    set -e;                                                                                                                                                                                                   \
        cd deploy/tektoncd/pipeline;                                                                                                                                                                              \
		rm templates/*.yaml;                                                                                                                                                                                      \
		mkdir temp;                                                                                                                                                                                               \
		cp release-v$(TEKTON_VERSION).yaml release.yaml;                                                                                                                                                          \
		kustomize build > temp/pipeline.yaml;                                                                                                                                                                     \
		cd temp;                                                                                                                                                                                                  \
        yq -s 'select(.) | .kind | downcase + $$index' pipeline.yaml;                                                                                                                                             \
		yq 'select(.) | .kind | downcase' pipeline.yaml | grep -v '\-\-\-' | uniq | xargs -I{} bash -c "cat {}[0-9]*.yml > ../templates/{}s.yaml";                                                                \
		cd ..;                                                                                                                                                                                                    \
		printf "apiVersion: v2\nappVersion: $(TEKTON_VERSION)\ndescription: A Helm chart for Tekton CD v$(TEKTON_VERSION).\nname: tektoncd-pipeline-operator\ntype: application\nversion: 0.1.0\n" > Chart.yaml;  \
		rm -rf temp release.yaml;                                                                                                                                                                                 \
	)

deploy-tektoncd-pipeline: tektoncd-pipeline-helm
	helm upgrade tektoncd ./deploy/tektoncd/pipeline \
		--install \
		--namespace $(TEKTON_NS) \
		--create-namespace

deploy/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml:
    wget "https://github.com/tektoncd/dashboard/releases/download/v$(TEKTON_DASHBOARD_VERSION)/tekton-dashboard-release.yaml" -O $@

tektoncd-dashboard-helm: deploy/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml
	$(shell                                                                                                                                                                                                                             \
	    set -e;                                                                                                                                                                                                                         \
        cd deploy/tektoncd/dashboard;                                                                                                                                                                                                   \
		find templates -name '*.yaml' -type f -not -name 'ingress.yaml' -delete;                                                                                                                                                        \
		mkdir -p temp;                                                                                                                                                                                                                  \
		cp release-v$(TEKTON_DASHBOARD_VERSION).yaml release.yaml;                                                                                                                                                                      \
		kustomize build > temp/release.yaml;                                                                                                                                                                                            \
		cd temp;                                                                                                                                                                                                                        \
		yq -s 'select(.) | .kind | downcase + $$index' release.yaml;                                                                                                                                                                    \
		yq 'select(.) | .kind | downcase' ../release-v$(TEKTON_DASHBOARD_VERSION).yaml | grep -v '\-\-\-' | uniq | xargs -I{} bash -c "cat {}[0-9]*.yml > ../templates/{}s.yaml";                                                       \
		cd ..;                                                                                                                                                                                                                          \
		printf "apiVersion: v2\nappVersion: $(TEKTON_DASHBOARD_VERSION)\ndescription: A Helm chart for Tekton CD dashboard v$(TEKTON_DASHBOARD_VERSION).\nname: tektoncd-dashboard\ntype: application\nversion: 0.1.0\n" > Chart.yaml;  \
		rm -rf temp release.yaml;                                                                                                                                                                                                       \
	)

deploy-tektoncd-dashboard: tektoncd-dashboard-helm
	helm upgrade tektoncd-dashboard ./deploy/tektoncd/dashboard \
		--install \
		--values ./deploy/tektoncd/dashboard/values.yaml \
		--namespace $(TEKTON_NS)

dev-deploy: deploy-nginx deploy-arangodb deploy-elasticsearch deploy-mongodb deploy-pg deploy-tektoncd-pipeline deploy-tektoncd-dashboard
