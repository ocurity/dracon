.PHONY: container publish

CONTAINER_REPO=
DRACON_VERSION=
SOURCE_CODE_REPO=
PRODUCER_AGGREGATOR_BASE_IMAGE=$(shell test -e .custom_image && cat .custom_image || echo "scratch")

DOCKER=docker

container:
	$(eval workdir:=$(shell mktemp -d /tmp/pdf.XXXXXX))
	mkdir -p ${workdir}/components/consumers && \
	cp -r ../../../bin/components/consumers/pdf ${workdir}/components/consumers && \
	cp default.html  ${workdir}/components/consumers/pdf && \
	$(DOCKER) build --tag $(CONTAINER_REPO)/components/consumers/pdf:$(DRACON_VERSION) \
					--file Dockerfile \
					$$([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
					${workdir} 1>&2 && \
	rm -rf ${workdir}

publish:
	$(DOCKER) push $(CONTAINER_REPO)/components/consumers/pdf:$(DRACON_VERSION) 1>&2
