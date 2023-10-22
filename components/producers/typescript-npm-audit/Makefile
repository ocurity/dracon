.PHONY build-extras

DOCKER_REPO=
DRACON_VERSION=

build-extras:
	docker build -t $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) -f tool/Dockerfile-tool ./tool/
