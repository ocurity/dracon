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
        "main_test.go",
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
