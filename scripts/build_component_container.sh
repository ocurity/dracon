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

if make -C "${executable_src_path}" --no-print-directory --dry-run container >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet container DOCKER_REPO="${DOCKER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
else
    dockerfile_template="
        FROM scratch                                    \n
        COPY ${executable_path} /app/${executable_path} \n
        ENTRYPOINT ["/app/${executable_path}"]          \n
    "
    dockerfile_path=$(mktemp)
    printf "${dockerfile_template}" > "${dockerfile_path}"
    docker build -t "${DOCKER_REPO}/${executable_src_path}:${DRACON_VERSION}" -f "${dockerfile_path}" ./bin
fi

if make -C "${executable_src_path}" --no-print-directory --dry-run extras >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet package-extras DOCKER_REPO="${DOCKER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
fi
