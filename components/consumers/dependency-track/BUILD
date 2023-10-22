subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "dependency-track",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
