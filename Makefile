component_binariess=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/bin')
component_containers=$(shell find ./components -name main.go | xargs -I'{}' sh -c 'echo $$(dirname {})/docker')
component_containers_publish=$(component_containers:docker=publish)
latest_tag=$(shell git tag --list --sort="-version:refname" | head -n 1)
commits_since_latest_tag=$(shell git log --oneline $(latest_tag)..HEAD | wc -l)

GO_TEST_PACKAGES=$(shell go list ./... | grep -v /vendor/)
CONTAINER_REPO=europe-west1-docker.pkg.dev/oc-dracon-saas/demo/ocurity/dracon
DRACON_VERSION=$(shell echo $(latest_tag)$$([ $(commits_since_latest_tag) -eq 0 ] || echo "-$$(git log -n 1 --pretty='format:%h')" )$$([ -z "$$(git status --porcelain=v1 2>/dev/null)" ] || echo "-dirty" ))

DOCKER=docker

export

########################################
############# BUILD TARGETS ############
########################################
.PHONY: components component-binaries build publish-component-containers

$(component_binariess):
	./scripts/build_component_binary.sh $@

component-binaries: $(component_binariess)

$(component_containers): %/docker: %/bin
	./scripts/build_component_container.sh $@

components: $(component_containers)

bin/cmd/kustomize-component-generator:
	@go build -o bin/cmd/kustomize-component-generator cmd/component-generator/main.go

build: bin/kustomize-component-generator components
	@echo "done building"

$(component_containers_publish): %/publish: %/docker
	./scripts/publish_component_container.sh $@

publish-component-containers: $(component_containers_publish)
