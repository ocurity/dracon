#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

mapfile -t images < \
    <(./pleasew query alltargets \
    --include buildkit-image \
    //... \
    | grep -v "^//third_party/")

util::info "Building ${#images[@]} image(s)"
./pleasew build "${images[@]}"

image_push_targets=()
for img in "${images[@]}"; do
    image_push_targets+=("${img}_push")
done
util::info "Pushing ${#image_push_targets[@]} image(s)"
./pleasew run parallel "${image_push_targets[@]}"
