.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=

DOCKER=docker

container:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/tagger:$(DRACON_VERSION) -f Dockerfile ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
