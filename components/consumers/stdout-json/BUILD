subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "stdout-json",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":stdout-json"],
)

dracon_component(
    name = "stdout-json",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
