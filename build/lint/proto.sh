#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

BUF="//third_party/binary/bufbuild/buf"

util::info "checking proto files"
"$BUF" lint . --path "./api/proto"

util::success "checked proto files"
