#!/bin/bash

# This script is used to update the open-source components in the local dracon 
# instance. It will build all components based on the latest local commit and 
# update the local dracon instance with the new components.

echo "Updating open-source components in local dracon instance..."
export CUSTOM_DRACON_VERSION=$(make print-DRACON_VERSION)
echo "DRACON_VERSION=$CUSTOM_DRACON_VERSION"
export CUSTOM_HELM_COMPONENT_PACKAGE_NAME="dracon-oss-components"

make publish-component-containers CONTAINER_REPO=localhost:5000/ocurity/dracon
bin/cmd/draconctl components package   --version ${CUSTOM_DRACON_VERSION}   \
  --chart-version ${CUSTOM_DRACON_VERSION}   \
  --name ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}   \
  ./components
helm upgrade ${CUSTOM_HELM_COMPONENT_PACKAGE_NAME} \
  ./${CUSTOM_HELM_COMPONENT_PACKAGE_NAME}-${CUSTOM_DRACON_VERSION}.tgz   \
  --install \
  --namespace dracon \
  --set container_registry=kind-registry:5000/ocurity/dracon

echo "Done! Bumped version to $CUSTOM_DRACON_VERSION"
