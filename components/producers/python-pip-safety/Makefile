.PHONY package-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

package-extras:
	$(docker) build -t $(DOCKER_REPO)/components/producers/python-pip-safety-tool:$(DRACON_VERSION) -f Dockerfile-tool .
