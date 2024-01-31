.PHONY: container

DOCKER_REPO=
DRACON_VERSION=
ESLINT_WRAPPER_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "node:lts")

DOCKER=docker

container:
	$(DOCKER) build --tag $(DOCKER_REPO)/components/producers/typescript-eslint/eslint-wrapper:$(DRACON_VERSION) \
					--file Dockerfile \
					--build-arg ESLINT_WRAPPER_BASE_IMAGE=$(ESLINT_WRAPPER_BASE_IMAGE) ../../../.. 1>&2
