#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

util::info "formatting BUILD files"
./pleasew fmt --write
util::success "formatted BUILD files"
