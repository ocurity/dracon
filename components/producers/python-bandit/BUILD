subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
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
