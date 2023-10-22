subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "semgrep",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
