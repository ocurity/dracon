.PHONY: build

DOCKER_REPO=
DRACON_VERSION=
EXECUTABLE_PATH=

DOCKER=docker

package:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/tagger:$(DRACON_VERSION) -f Dockerfile ../../../bin
