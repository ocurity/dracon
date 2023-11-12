.PHONY: build clean clean-protos lint fmt fmt-go fmt-proto tests go-tests release_notes kustomizations

proto_defs=$(shell find . -name "*.proto" -not -path "./vendor/*")
go_protos=$(proto_defs:.proto=.pb.go)
component_containers=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/docker')
component_tasks=$(shell find ./components -name task.yaml -not -path "./components/sources/git/*" -not -path "./components/enrichers/aggregator/*" -not -path "./components/producers/aggregator/*" -not -path "./components/base/*")
component_kustomizations=$(component_tasks:task.yaml=kustomization.yaml)

GO_TEST_PACKAGES=./...

DOCKER_REPO=ghcr.io/ocurity/dracon
TEKTON_VERSION=0.44.0
TEKTON_DASHBOARD_VERSION=0.29.2

latest_tag=$(shell git tag --list --sort="-version:refname" | head -n 1)
commits_since_latest_tag=$(shell git log --oneline $(latest_tag)..HEAD | wc -l)
DRACON_VERSION=$(shell echo $(latest_tag)$$([ $(commits_since_latest_tag) -eq 0 ] || echo "-$$(git log -n 1 --pretty='format:%h')" )$$([ -z "$$(git status --porcelain=v1 2>/dev/null)" ] || echo "-dirty" ))

PROTOC=protoc

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

$(component_containers):
	$(shell                                                                                                                                         \
        set -e;                                                                                                                                     \
        EXECUTABLE=$$(basename $$(dirname $@));                                                                                                     \
        echo "$@" | grep -Eq ^components/producers/.*$ && EXECUTABLE=$${EXECUTABLE}-parser || true;                                                 \
        EXECUTABLE_FOLDER=$$(dirname $@);                                                                                                           \
        EXECUTABLE_PATH=$$(dirname $@)/$${EXECUTABLE};                                                                                              \
        if $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --dry-run build >/dev/null 2>&1;                                                   \
        then                                                                                                                                        \
            $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --quiet build DOCKER_REPO=$(DOCKER_REPO) DRACON_VERSION=$(DRACON_VERSION);        \
        else                                                                                                                                        \
            DOCKERFILE_TEMPLATE="                                                                                                                   \
                FROM golang:1.21.3-bookworm as builder                                                                                           \n \
                COPY . /build                                                                                                                    \n \
                WORKDIR /build                                                                                                                   \n \
                RUN go build -o $${EXECUTABLE_PATH} ./$${EXECUTABLE_FOLDER}/main.go                                                              \n \
                FROM scratch                                                                                                                     \n \
                COPY --from=builder /build/$${EXECUTABLE_PATH} /bin/$${EXECUTABLE}                                                               \n \
                ENTRYPOINT ["/bin/$${EXECUTABLE}"]                                                                                               \n \
            ";                                                                                                                                      \
            DOCKERFILE=$$(mktemp);                                                                                                                  \
            printf "$${DOCKERFILE_TEMPLATE}" > $${DOCKERFILE};                                                                                      \
            docker build -t $(DOCKER_REPO)/$${EXECUTABLE_FOLDER}:$(DRACON_VERSION) -f $${DOCKERFILE} .;                                             \
        fi;                                                                                                                                         \
        if $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --dry-run build-extras >/dev/null 2>&1;                                            \
        then                                                                                                                                        \
            $(MAKE) -C $${EXECUTABLE_FOLDER} --no-print-directory --quiet build-extras DOCKER_REPO=$(DOCKER_REPO) DRACON_VERSION=$(DRACON_VERSION); \
        fi;                                                                                                                                         \
	)

components: components/base/openapi_schema.json $(component_containers)

third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json:
	wget "https://raw.githubusercontent.com/tektoncd/pipeline/v$(TEKTON_VERSION)/pkg/apis/pipeline/v1beta1/swagger.json" -O third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json

third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml:
	wget "https://storage.googleapis.com/tekton-releases/pipeline/previous/v$(TEKTON_VERSION)/release.yaml" -O third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml

api/openapi/tekton/openapi_schema.json: third_party/k8s/tektoncd/pipeline/swagger-v$(TEKTON_VERSION).json
	./scripts/generate_openapi_schema.sh $< $@

components/base/openapi_schema.json: api/openapi/tekton/openapi_schema.json
	cp $< $@

mirror_images:
	./scripts/mirror_images.sh

third_party/k8s/tektoncd/pipeline/pipeline.yaml: third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml
	cp third_party/k8s/tektoncd/pipeline/release-v$(TEKTON_VERSION).yaml third_party/k8s/tektoncd/pipeline/release.yaml
	kustomize build $(shell dirname $@) > $@
	rm third_party/k8s/tektoncd/pipeline/release.yaml

third_party/k8s/tektoncd/pipeline/Chart.yaml: third_party/k8s/tektoncd/pipeline/pipeline.yaml
	printf "apiVersion: v2\nappVersion: $(TEKTON_VERSION)\ndescription: A Helm chart for Kubernetes.\nname: pipeline\ntype: application\nversion: 0.1.0\n" > $@

third_party/k8s/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml:
    wget "https://github.com/tektoncd/dashboard/releases/download/v$(TEKTON_DASHBOARD_VERSION)/tekton-dashboard-release.yaml" -O third_party/k8s/tektoncd/dashboard/release-v$(TEKTON_DASHBOARD_VERSION).yaml

release_notes:
	git log --date=short --pretty='format:- %cd %s' -n 20

bin/kustomize-component-generator:
	go build -o bin/kustomize-component-generator build/tools/kustomize-component-generator/main.go

$(component_kustomizations): bin/kustomize-component-generator
	bin/kustomize-component-generator -task "$$(dirname $@)/task.yaml"

kustomizations: $(component_kustomizations)

print-%:
	@echo $($*)
