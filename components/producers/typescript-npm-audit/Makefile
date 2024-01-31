.PHONY: extras publish-extras

DOCKER_REPO=
DRACON_VERSION=
TYPESCRIPT_NPM_AUDIT_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "node:15-alpine3.12")

DOCKER=docker

extras:
	$(DOCKER) build --tag $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) \
					--file tool/Dockerfile-tool \
					--build-arg TYPESCRIPT_NPM_AUDIT_BASE_IMAGE=$(TYPESCRIPT_NPM_AUDIT_BASE_IMAGE) ./tool/ 1>&2

publish-extras:
	$(DOCKER) push $(DOCKER_REPO)/components/producers/typescript-npm-audit-tool:$(DRACON_VERSION) 1>&2
