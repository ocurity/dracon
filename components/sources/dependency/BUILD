subinclude(
    "///k8s//build/defs:k8s",
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "source-dependency",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
