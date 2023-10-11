#!/usr/bin/env bash
# This script performs linting for Dracon.
set -Eeuo pipefail

path_dirs=(
    "$(dirname "//third_party/binary/dominikh/go-tools:staticcheck")"
    "$(dirname "//third_party/binary/mgechev/revive:revive")"
    "$(dirname "//third_party/binary/securego/gosec:gosec")"
    "$(dirname "//third_party/go:mvdan.cc_gofumpt")"
    "$(echo "//third_party/binary/bufbuild/buf:buf")/bin"
    "$(echo "//third_party/go:toolchain")/bin"
    "$(dirname "//third_party/binary/reviewdog/reviewdog")"
)

for pd in "${path_dirs[@]}"; do
    export PATH="$pd:$PATH"
done

export GO_PACKAGES="$(
    ./pleasew query alltargets --include go \
    | grep -v "//third_party/" \
    | cut -f3- -d/ \
    | cut -f1 -d: \
    | awk '{ print "./" $0 }' \
    | xargs
)"

./pleasew -p -v=2 generate //...
./pleasew -p -v=2 build --include codegen //...

cmd=(
    "reviewdog"
    "-fail-on-error"
    "$@"
)

if [[ "${CI:-}" == "true" ]]; then
    cmd+=(
        "-reporter=github-pr-review"
        "-tee"
    )
elif [[ "$@" != *"-filter-mode=nofilter"* ]]; then
    cmd+=(
        "-diff="git diff origin/main""
    )
fi
set -x
"${cmd[@]}"
