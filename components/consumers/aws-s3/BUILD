subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "aws-s3",
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
    srcs = [":aws-s3"],
)

dracon_component(
    name = "aws-s3",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
