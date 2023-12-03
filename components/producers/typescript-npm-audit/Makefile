.PHONY package-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

package-extras:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) -f tool/Dockerfile-tool ./tool/
