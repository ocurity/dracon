#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

BUF="//third_party/binary/bufbuild/buf"

util::info "formatting proto files"
"$BUF/bin/buf" format -w "./api/proto"

util::success "formatted proto files"
