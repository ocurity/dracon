subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
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
