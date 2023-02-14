subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "python-pip-safety-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/python-pip-safety/types",
    ],
)

go_test(
    name = "python-pip-safety-test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/python-pip-safety/types",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":python-pip-safety-parser"],
)

dracon_component(
    name = "python-pip-safety",
    images = [
        ":image",
        ":tool",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)

buildkit_image(
    name = "tool",
    dockerfile = "Dockerfile-tool",
)
