subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "aws-s3",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
        "//third_party/go/github.com/aws/aws-sdk-go",
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
