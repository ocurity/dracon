subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "python-bandit-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
    ],
)

go_test(
    name = "python-bandit_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/context",
        "//pkg/testutil",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":python-bandit-parser"],
    visibility = [
        "//examples/...",
    ],
)

buildkit_image(
    name = "tool",
    dockerfile = "Dockerfile-tool",
)

dracon_component(
    name = "python-bandit",
    images = [
        ":image",
        ":tool",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
