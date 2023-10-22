subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "stdout-json",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
