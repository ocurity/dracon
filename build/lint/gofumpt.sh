#!/usr/bin/env bash
set -Eeuo pipefail

GOFUMPT="//third_party/go:mvdan.cc_gofumpt"

mapfile -t go_dirs < <(./pleasew query alltargets \
    --include=go \
    | grep -v third_party \
    | cut -f1 -d":" \
    | cut -c 3- \
    | sort -u
)

if [ -v go_dirs ]; then
    "$GOFUMPT" -d -e "${go_dirs[@]}"
fi
