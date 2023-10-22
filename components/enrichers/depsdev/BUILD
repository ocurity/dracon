subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

dracon_component(
    name = "depsdev",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
