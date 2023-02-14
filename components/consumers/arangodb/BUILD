subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "arangodb",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//third_party/go/github.com/arangodb/go-driver",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":arangodb"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "arangodb",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
