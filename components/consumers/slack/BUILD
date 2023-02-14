subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "slack",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//components/consumers/slack/utils",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":slack"],
)

dracon_component(
    name = "slack",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
