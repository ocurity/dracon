subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "pdf",
    srcs = [
        "main.go",
        "default.html",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
        "//pkg/templating",
        "//third_party/go/google.golang.org/protobuf",
        "//third_party/go/github.com/aws/aws-sdk-go",
        "//third_party/go/github.com/playwright-community/playwright-go",
    ],
)

# buildkit_distroless_image(
#     name = "image",
#     srcs = [":pdf"],
#     visibility = [
#         "//examples/...",
#     ],
# )

buildkit_image(
    name = "image",
    srcs = [
        ":pdf",
    ],
    dockerfile = "Dockerfile",
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "pdf",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
