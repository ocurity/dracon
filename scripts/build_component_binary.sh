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
echo "${1}" | grep -Eq ^components/enrichers/.*$ && executable="${executable}-enricher" || true
echo "${1}" | grep -Eq ^components/consumers/.*$ && executable="${executable}-consumer" || true

executable_src_path=$(dirname "${1}")
executable_path=$(dirname $(dirname "${1}"))/"${executable}"

echo "building bin/${executable_path}/${executable}" > /dev/stderr

go build -o "bin/${executable_src_path}/${executable}" "./${executable_src_path}/main.go"
