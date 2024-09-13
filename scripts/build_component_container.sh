#!/bin/bash

set -e;

source ./scripts/util.sh
containers_path=./containers

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

if make -C "${executable_src_path}" --no-print-directory --dry-run container >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet container CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
else
    dockerfile_template=""
    if [[ -n "${BASE_IMAGE+x}" ]]
    then
        echo "Using base image: ${BASE_IMAGE}"
        dockerfile_template=$(cat "${containers_path}/Dockerfile.base.image")
    else
        BASE_IMAGE=''
        dockerfile_template=$(cat "${containers_path}/Dockerfile.base")
    fi

    dockerfile_path=$(mktemp)
    printf "${dockerfile_template}" > "${dockerfile_path}"
    docker build \
        --platform "${GOOS}/${GOARCH}" \
        --build-arg EXECUTABLE_PATH=${executable_path} \
        --build-arg BASE_IMAGE=${BASE_IMAGE} \
        -t "${CONTAINER_REPO}/${executable_src_path}:${DRACON_VERSION}" \
        $([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
        -f "${dockerfile_path}" ./bin
fi

if make -C "${executable_src_path}" --no-print-directory --dry-run extras >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet extras CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
fi
