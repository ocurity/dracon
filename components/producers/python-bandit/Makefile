.PHONY package-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

package-extras:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/python-bandit-tool:$(DRACON_VERSION) -f Dockerfile-tool .
