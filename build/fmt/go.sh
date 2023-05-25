#!/usr/bin/env bash

set -Eeuo pipefail

source "//build/util"

GOFUMPT="//third_party/go:mvdan.cc_gofumpt"
GO_ROOT="//third_party/go:toolchain"

util::info "formatting go files"

mapfile -t go_dirs < <(./pleasew query alltargets \
    --include=go \
    | grep -v third_party \
    | cut -f1 -d":" \
    | cut -c 3- \
    | sort -u
)

if [ -v go_dirs ]; then
    "$GOFUMPT" -w "${go_dirs[@]}"
fi

util::success "formatted go files"

util::info "tidying go module"

GO111MODULE=on "${GO_ROOT}/bin/go" mod tidy
util::success "tidied go module"
