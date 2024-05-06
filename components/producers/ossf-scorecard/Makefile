.PHONY: extras publish-extras

CONTAINER_REPO=
DRACON_VERSION=
OSSF_SCORECARD_SAFETY_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "alpine:latest")

DOCKER=docker

extras:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					--build-arg OSSF_SCORECARD_SAFETY_BASE_IMAGE=$(OSSF_SCORECARD_SAFETY_BASE_IMAGE) ../../.. 1>&2

publish-extras:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) 1>&2
