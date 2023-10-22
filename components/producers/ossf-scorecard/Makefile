.PHONY build-extras

DOCKER_REPO=
DRACON_VERSION=

build-extras:
	docker build -t $(DOCKER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) -f Dockerfile ../../..
