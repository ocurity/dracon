.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=
PRODUCER_AGGREGATOR_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "scratch")

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/consumers/dependency-track:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					 ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/consumers/dependency-track:$(DRACON_VERSION) 1>&2
