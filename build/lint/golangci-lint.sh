#!/usr/bin/env bash
set -Eeuo pipefail

GOLANGCI_LINT="//build/util:golangci-lint"

mapfile -t go_dirs < <(./pleasew query alltargets \
    --include=go \
    | grep -v third_party \
    | cut -f1 -d":" \
    | cut -c 3- \
    | sort -u
)

if [ -v go_dirs ]; then
    "$GOLANGCI_LINT" run "${go_dirs[@]}"
fi
