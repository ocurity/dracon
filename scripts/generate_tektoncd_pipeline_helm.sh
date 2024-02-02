#!/bin/bash

set -e

source ./scripts/util.sh

if [ "$#" -eq 0 ]
then
    util::error "No directory provided to build"
    exit 1
fi

TEKTON_VERSION="${1}"

cd deploy/tektoncd/pipeline
find templates -name '*.yaml' -type f -delete
mkdir -p temp

cp release-v${TEKTON_VERSION}.yaml release.yaml
kustomize build > temp/pipeline.yaml

cd temp
yq -s 'select(.) | .kind | downcase + $index' pipeline.yaml
yq 'select(.) | .kind | downcase' pipeline.yaml | grep -v '\-\-\-' | uniq | xargs -I{} bash -c "cat {}[0-9]*.yml > ../templates/{}s.yaml"
cd ..

cat <<EOF> Chart.yaml
apiVersion: v2
appVersion: ${TEKTON_VERSION}
description: A Helm chart for Tekton CD v${TEKTON_VERSION}.
name: tektoncd-pipeline-operator
type: application
version: 0.1.0
EOF

rm -rf temp release.yaml
