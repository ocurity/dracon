.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=
PRODUCER_AGGREGATOR_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "scratch")
# https://docs.docker.com/build/building/multi-platform/
# Make sure to always build containers using AMD64 but allow to be overridden by users if need for cross-os compatibility.
CONTAINER_OS_ARCH=linux/amd64
# Allow to independently customise go OS and ARCH flags.
# Defaulting to linux/amd64 as per CONTAINER_OS_ARCH.
GOOS=linux
GOARCH=amd64

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/consumers/dependency-track:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg GOOS=$(GOOS) \
					--build-arg GOARCH=$(GOARCH) \
					../../../bin 1>&2 \
					--platform=$(CONTAINER_OS_ARCH)

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/consumers/dependency-track:$(DRACON_VERSION) 1>&2
