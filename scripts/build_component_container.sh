#!/bin/bash

set -e;

source ./scripts/util.sh

# Sanity check for not arguments being passed.
if [ "$#" -eq 0 ]
then
    util::error "No arguments provided to build. Expected two."
    exit 1
fi

if [ -z "$1" ]
then
    util::error "No directory argument provided to build."
    exit 1
fi

if [ -z "$2" ]
then
    util::error "No build architecture argument provided to build."
    exit 1
fi

dir_name="$1"
build_architecture="$2"

executable=$(basename $(dirname ${dir_name}))

echo ${dir_name} | grep -Eq ^components/producers/.*$ && executable="${executable}-parser" || true
echo ${dir_name} | grep -Eq ^components/enrichers/.*$ && executable="${executable}" || true
echo ${dir_name} | grep -Eq ^components/consumers/.*$ && executable="${executable}" || true

executable_src_path=$(dirname $dir_name)
executable_path=$(dirname $dir_name)/${build_architecture}/"${executable}"

if make -C "${executable_src_path}" --no-print-directory --dry-run container >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet container CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
else
    dockerfile_template="
        FROM ${BASE_IMAGE:-scratch}                \n
        COPY ${executable_path} /app/${executable} \n
        ENTRYPOINT ["/app/${executable}"]          \n
    "
    dockerfile_path=$(mktemp)
    printf "${dockerfile_template}" > "${dockerfile_path}"
    docker build -t "${CONTAINER_REPO}/${executable_src_path}:${DRACON_VERSION}" \
        $([ "${SOURCE_CODE_REPO}" != "" ] && echo "--label=org.opencontainers.image.source=${SOURCE_CODE_REPO}" ) \
        -f "${dockerfile_path}" ./bin \
        --platform=${build_architecture}
fi

if make -C "${executable_src_path}" --no-print-directory --dry-run extras >/dev/null 2>&1
then
    make -C "${executable_src_path}" --no-print-directory --quiet extras CONTAINER_REPO="${CONTAINER_REPO}" DRACON_VERSION="${DRACON_VERSION}"
fi
