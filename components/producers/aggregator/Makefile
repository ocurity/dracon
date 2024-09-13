.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=

DOCKER=docker

container:
	$(DOCKER) --debug build --tag $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) \
					--build-arg EXECUTABLE_SRC_PATH=components/producers/aggregator/aggregator-parser \
        			--build-arg EXECUTABLE_TARGET_PATH=components/producers/aggregator/tagger \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--file ${BASE_IMAGE_PATH} ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
