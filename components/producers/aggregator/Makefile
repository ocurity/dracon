.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg BASE_IMAGE=$(BASE_IMAGE) ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
