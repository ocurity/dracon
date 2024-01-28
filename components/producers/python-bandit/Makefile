.PHONY: extras publish-extras

CONTAINER_REPO=
DRACON_VERSION=

DOCKER=docker

extras:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/python-bandit-tool:$(DRACON_VERSION) -f Dockerfile-tool . 1>&2

publish-extras:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/python-bandit-tool:$(DRACON_VERSION) 1>&2
