.PHONY: container

CONTAINER_REPO=
DRACON_VERSION=

# https://docs.docker.com/build/building/multi-platform/
# Make sure to always build containers using AMD64 but allow to be overridden by users if need for cross-os compatibility.
GOOS=linux
GOARCH=amd64

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/typescript-eslint/eslint-wrapper:$(DRACON_VERSION) \
					--platform $(GOOS)/$(GOARCH) \
					--build-arg GOOS=$(GOOS) \
					--build-arg GOARCH=$(GOARCH) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg ESLINT_WRAPPER_BASE_IMAGE=$(ESLINT_WRAPPER_BASE_IMAGE) ../../../.. 1>&2
