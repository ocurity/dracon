subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "jira",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//components/consumers/jira/utils",
        "//pkg/jira/config",
        "//pkg/jira/jira",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":jira"],
)

dracon_component(
    name = "jira",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
