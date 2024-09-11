.PHONY: container

CONTAINER_REPO=
DRACON_VERSION=
ESLINT_WRAPPER_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "node:lts")
# https://docs.docker.com/build/building/multi-platform/
# Make sure to always build containers using AMD64 but allow to be overridden by users if need for cross-os compatibility.
CONTAINER_OS_ARCH=linux/amd64
# Allow to independently customise go OS and ARCH flags.
# Defaulting to linux/amd64 as per CONTAINER_OS_ARCH.
GOOS=linux
GOARCH=amd64

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/typescript-eslint/eslint-wrapper:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg GOOS=$(GOOS) \
					--build-arg GOARCH=$(GOARCH) \
					--build-arg ESLINT_WRAPPER_BASE_IMAGE=$(ESLINT_WRAPPER_BASE_IMAGE) ../../../.. 1>&2 \
					--platform=$(CONTAINER_OS_ARCH)
