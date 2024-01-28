.PHONY: container

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

container:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/typescripit-eslint/eslint-wrapper:$(DRACON_VERSION) -f Dockerfile ../../../.. 1>&2
