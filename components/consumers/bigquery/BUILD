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
        "//third_party/go/cloud.google.com/go/bigquery",
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
