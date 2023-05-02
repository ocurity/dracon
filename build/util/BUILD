subinclude("//build/defs:sh")

sh_script(
    name = "util",
    main = "util.sh",
    visibility = ["//build/..."],
    deps = [
        "//third_party/sh:ansi",
    ],
)

sh_cmd(
    name = "lint",
    cmd = """
    ./pleasew -p -v=2 generate
    cmd=(
        "$(out_exe //third_party/binary/reviewdog/reviewdog)"
        -fail-on-error
    )
    if [[ "\\\\$CI" == "true" ]]; then
        cmd+=(
            "-reporter=github-pr-review"
        )
    elif [[ "\\\\$@" != *"-filter-mode=nofilter"* ]]; then
        cmd+=(
            "-diff="git diff origin/main""
        )
    fi
    set -x
    "\\\\${cmd[@]}" \\\\$@
    """,
    data = ["//third_party/binary/reviewdog/reviewdog"],
    shell = "/usr/bin/env bash",
)

sh_cmd(
    name = "fmt_all",
    cmd = """
    ./pleasew run sequential \\\\
        //build/fmt:please \\\\
        //build/fmt:go \\\\
        //build/fmt:proto
    """,
)

sh_cmd(
    name = "dev_deploy",
    cmd = """
set -Eeuo pipefail
labels="//build/k8s/k3d:dracon,support"
if [ -n "\\\\${1:-}" ]; then
    labels="\\\\$labels,\\\\$1"
fi
set -x
./pleasew run ///k8s//tools/deployer:deployer -- pipeline_deploy \\
    --labels="\\\\$labels" \\
    --ordered_labels="webhook;operator"
    """,
    shell = "/usr/bin/env bash",
)

sh_cmd(
    name = "golangci-lint",
    cmd = """
    export PATH="\\\\$(./pleasew query reporoot)/$(out_location //third_party/go:toolchain)/bin:\\\\$PATH"
    go mod tidy
    $(out_exe //third_party/binary/golangci/golangci-lint) "\\\\$@"
    """,
    data = [
        "//third_party/binary/golangci/golangci-lint",
        "//third_party/go:toolchain",
    ],
    visibility = ["//build/lint:golangci-lint"],
)
