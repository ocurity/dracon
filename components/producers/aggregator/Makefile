.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=

# https://docs.docker.com/build/building/multi-platform/
# Make sure to always build containers using AMD64 but allow to be overridden by users if need for cross-os compatibility.
GOOS=linux
GOARCH=amd64

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) \
					--platform $(GOOS)/$(GOARCH) \
					--build-arg GOOS=$(GOOS) \
					--build-arg GOARCH=$(GOARCH) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg BASE_IMAGE=$(BASE_IMAGE) ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
