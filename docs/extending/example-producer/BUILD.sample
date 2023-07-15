subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "example-tool-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
    ],
)

go_test(
    name = "example-tool-parser_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":example-tool-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "example-tool",
    images = [
        ":image",
        "//third_party/docker/example/tool",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
