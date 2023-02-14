#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

util::info "checking BUILD files"
if ! ./pleasew fmt --quiet; then
  util::error "BUILD files incorrectly formatted. Please run:
  $ ./pleasew run //scripts/fmt:plz"
  exit 1
fi
util::success "checked BUILD files"

if ! [[ $(git status --porcelain) ]]; then
    util::info "checking Generated files"
    ./pleasew generate
    if [[ $(git status --porcelain) ]]; then
    util::error "./pleasew generate made changes to files. Please commit them."
    exit 1
    fi
    util::success "checked Generated files"
fi
