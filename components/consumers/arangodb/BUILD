subinclude("///go//build_defs:go")

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
        "///third_party/go/github.com_arangodb_go-driver//:go-driver",
        "///third_party/go/github.com_arangodb_go-driver//http",
        "//api/proto/v1",
        "//components/consumers",
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
