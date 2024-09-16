#!/bin/bash

set -e

source ./scripts/util.sh

if [ "$#" -eq 0 ]
then
    util::error "No directory provided to build"
    exit 1
fi

executable=$(basename $(dirname "${1})"))

echo "${1}" | grep -Eq ^components/producers/.*$ && executable="${executable}-parser" || true
echo "${1}" | grep -Eq ^components/enrichers/.*$ && executable="${executable}" || true
echo "${1}" | grep -Eq ^components/consumers/.*$ && executable="${executable}" || true

executable_src_path=$(dirname "${1}")
executable_path=$(dirname $(dirname "${1}"))/"${executable}"

# Customised binary per OS/ARCH.
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
out_bin_path="bin/${executable_src_path}/${GOOS}/${GOARCH}/${executable}"

echo "building $out_bin_path" > /dev/stderr

CGO_ENABLED=0 go build -o $out_bin_path "./${executable_src_path}/main.go"
