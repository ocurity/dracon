.PHONY build-extras

DOCKER_REPO=
DRACON_VERSION=

build-extras:
	docker build -t $(DOCKER_REPO)/components/producers/python-pip-safety-tool:$(DRACON_VERSION) -f Dockerfile-tool .
