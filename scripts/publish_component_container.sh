#!/bin/bash

set -e

source ./scripts/util.sh

if [ "$#" -eq 0 ]
then
    util::error "No directory provided to build"
    exit 1
fi

if make -C $(dirname "${1}") --no-print-directory --dry-run publish >/dev/null 2>&1
then
	make -C $(dirname "${1}") --no-print-directory --quiet publish DOCKER_REPO="${DOCKER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
else
	docker push "${DOCKER_REPO}/$(dirname ${1}):${DRACON_VERSION}" 1>&2
fi

if make -C $(dirname "${1}") --no-print-directory --dry-run publish-extras >/dev/null 2>&1
then
	make -C $(dirname "${1}") --no-print-directory --quiet publish-extras DOCKER_REPO="${DOCKER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
fi
