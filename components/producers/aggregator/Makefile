.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
PRODUCER_AGGREGATOR_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "scratch")

DOCKER=docker

container:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) \
					--file Dockerfile \
					--build-arg PRODUCER_AGGREGATOR_BASE_IMAGE=$(PRODUCER_AGGREGATOR_BASE_IMAGE) ../../../bin 1>&2

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/tagger:$(DRACON_VERSION) 1>&2
