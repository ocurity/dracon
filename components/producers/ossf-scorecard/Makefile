.PHONY package-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

package-extras:
	$(docker) build -t $(DOCKER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) -f Dockerfile ../../..
