subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "mongodb",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//third_party/go/go.mongodb.org/mongo-driver",
        "//third_party/go/github.com/avast/retry-go/v4",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":mongodb"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "mongodb",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
