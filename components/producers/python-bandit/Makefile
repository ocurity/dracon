.PHONY: extras publish-extras

CONTAINER_REPO=
DRACON_VERSION=
PYTHON_BANDIT_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "python:3.10-slim")

DOCKER=docker

extras:
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/producers/python-bandit-tool:$(DRACON_VERSION) \
					--file Dockerfile-tool \
					--build-arg PYTHON_BANDIT_BASE_IMAGE=$(PYTHON_BANDIT_BASE_IMAGE) . 1>&2

publish-extras:
	$(DOCKER) push $(CONTAINER_REPO)/components/producers/python-bandit-tool:$(DRACON_VERSION) 1>&2
