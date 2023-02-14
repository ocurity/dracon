#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

SHELLCHECK="//third_party/binary/koalaman/shellcheck"

util::info "checking shell files"
mapfile -t sh_dirs < <(./pleasew query alltargets \
    --include sh \
    | cut -f1 -d":" \
    | cut -c 3- \
    | sort -u
)

export SHELLCHECK_OPTS="-e SC1091" # (info): Not following: //third_party/sh:shflags: openBinaryFile: does not exist (No such file or directory)

for dir in "${sh_dirs[@]}"; do
    mapfile -t files < <(find "${dir}/" -type f -name '*.sh')
    if ! "$SHELLCHECK" --external-sources "${files[@]}"; then
        util::error "shell files failed check"
        exit 1
    fi
done

# always use `./pleasew` instead of `plz`
if grep -ri "plz " "${sh_dirs[@]}" | grep -v "lint/sh"; then
    util::error "found usage of 'plz', please replaces it with './pleasew'"
    exit 1
fi

util::success "checked shell files"
