#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

ACTIONLINT="//third_party/binary/rhysd/actionlint"

util::info "checking GitHub Actions Workflows"
"$ACTIONLINT" -shellcheck= -pyflakes= -ignore 'machinetype' -ignore 'machineconfig'

util::success "checked GitHub Actions Workflows"
