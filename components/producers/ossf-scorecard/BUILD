subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

dracon_component(
    name = "ossf-scorecard",
    images = [
        ":image",
        ":scorecard-dracon",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
