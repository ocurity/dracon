subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "bigquery",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
        "//third_party/go/cloud.google.com/go/bigquery",
        "//third_party/go/golang.org/x/oauth2",
        "//third_party/go/google.golang.org/api",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":bigquery"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "bigquery",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
