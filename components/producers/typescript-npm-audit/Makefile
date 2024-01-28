.PHONY: extras publish-extras

DOCKER_REPO=
DRACON_VERSION=

DOCKER=docker

extras:
	$(DOCKER) build -t $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) -f tool/Dockerfile-tool ./tool/ 1>&2

publish-extras:
	$(DOCKER) push $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) 1>&2
