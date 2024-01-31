.PHONY: extras publish-extras

DOCKER_REPO=
DRACON_VERSION=
OSSF_SCORECARD_SAFETY_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "alpine:latest")

DOCKER=docker

extras:
	$(DOCKER) build --tag $(DOCKER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) \
					--file Dockerfile \
					--build-arg OSSF_SCORECARD_SAFETY_BASE_IMAGE=$(OSSF_SCORECARD_SAFETY_BASE_IMAGE) ../../.. 1>&2

publish-extras:
	$(DOCKER) push $(DOCKER_REPO)/components/producers/ossf-scorecard:$(DRACON_VERSION) 1>&2
