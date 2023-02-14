subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "npm-audit-parser",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-npm-audit/types",
        "//components/producers/typescript-npm-audit/types/npmfullaudit",
        "//components/producers/typescript-npm-audit/types/npmquickaudit",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":npm-audit-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "npm-audit",
    images = [
        ":image",
        "//third_party/docker:node",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
