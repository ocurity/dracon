subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)
# this producer covers yarn audit https://classic.yarnpkg.com/lang/en/docs/cli/audit/

go_binary(
    name = "yarn-audit-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-yarn-audit/types",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":yarn-audit-parser"],
)

dracon_component(
    name = "typescript-yarn-audit",
    images = [
        ":image",
        "//third_party/docker:node",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
