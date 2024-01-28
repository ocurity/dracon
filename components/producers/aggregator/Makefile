.PHONY: container publish

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

container:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/tagger:$(DRACON_VERSION) -f Dockerfile ../../../bin 1>&2

publish:
	$(DOCKER) push $(DOCKER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
