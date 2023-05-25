subinclude("//build/defs:sh")

sh_script(
    name = "util",
    main = "util.sh",
    visibility = ["//build/..."],
    deps = [
        "//third_party/sh:ansi",
    ],
)

sh_script(
    name = "lint",
    find_repo_root = True,
    main = "lint.sh",
    deps = [
        "//third_party/binary/bufbuild/buf",
        "//third_party/binary/dominikh/go-tools:staticcheck",
        "//third_party/binary/mgechev/revive",
        "//third_party/binary/reviewdog/reviewdog",
        "//third_party/binary/securego/gosec",
        "//third_party/go:mvdan.cc_gofumpt",
        "//third_party/go:toolchain",
    ],
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
