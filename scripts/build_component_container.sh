#!/bin/bash

set -e;

source ./scripts/util.sh

if [ "$#" -eq 0 ]
then
    util::error "No directory provided to build"
    exit 1
fi

executable=$(basename $(dirname "${1}"))

echo "${1}" | grep -Eq ^components/producers/.*$ && executable="${executable}-parser" || true
echo "${1}" | grep -Eq ^components/enrichers/.*$ && executable="${executable}" || true
echo "${1}" | grep -Eq ^components/consumers/.*$ && executable="${executable}" || true

executable_src_path=$(dirname "${1}")
executable_path=$(dirname "${1}")/"${executable}"

BASE_IMAGE_PATH=$(realpath ${BASE_IMAGE_PATH:-./containers/Dockerfile.base})

if make -C "${executable_src_path}" --no-print-directory --dry-run container >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet container BASE_IMAGE_PATH="${BASE_IMAGE_PATH}" CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
else
    docker build \
        --build-arg EXECUTABLE_SRC_PATH=${executable_path} \
        --build-arg EXECUTABLE_TARGET_PATH=${executable_path} \
        -t "${CONTAINER_REPO}/${executable_src_path}:${DRACON_VERSION}" \
        $([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
        -f "${BASE_IMAGE_PATH}" ./bin
fi

if make -C "${executable_src_path}" --no-print-directory --dry-run extras >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet extras CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
fi
