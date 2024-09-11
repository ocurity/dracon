.PHONY: container

CONTAINER_REPO=
DRACON_VERSION=
BUILD_ARCHITECTURE=

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/typescript-eslint/eslint-wrapper:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg ESLINT_WRAPPER_BASE_IMAGE=$(ESLINT_WRAPPER_BASE_IMAGE) \
					--build-arg BUILD_ARCHITECTURE=$(BUILD_ARCHITECTURE) \
					../../../.. 1>&2
