.PHONY: extras publish-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

extras:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/python-pip-safety-tool:$(DRACON_VERSION) -f Dockerfile-tool . 1>&2

publish-extras:
	$(DOCKER) push $(DOCKER_REPO)/components/producers/python-pip-safety-tool:$(DRACON_VERSION) 1>&2
