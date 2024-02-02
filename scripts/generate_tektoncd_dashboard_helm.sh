#!/bin/bash

set -e

source ./scripts/util.sh

if [ "$#" -eq 0 ]
then
    util::error "No directory provided to build"
    exit 1
fi

TEKTON_DASHBOARD_VERSION="${1}"

cd deploy/tektoncd/dashboard
find templates -name '*.yaml' -type f -not -name 'ingress.yaml' -delete
mkdir -p temp

cp release-v${TEKTON_DASHBOARD_VERSION}.yaml release.yaml
kustomize build > temp/release.yaml

cd temp
yq -s 'select(.) | .kind | downcase + $index' release.yaml
yq 'select(.) | .kind | downcase' ../release-v${TEKTON_DASHBOARD_VERSION}.yaml | grep -v '\-\-\-' | uniq | xargs -I{} bash -c "cat {}[0-9]*.yml > ../templates/{}s.yaml"
cd ..                                                                                                                                                                                                                          \

cat <<EOF> Chart.yaml
apiVersion: v2
appVersion: ${TEKTON_DASHBOARD_VERSION}
description: A Helm chart for Tekton CD dashboard v${TEKTON_DASHBOARD_VERSION}.
name: tektoncd-dashboard
type: application
version: 0.1.0
EOF

rm -rf temp release.yaml