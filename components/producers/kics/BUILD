subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "kics",
    images = [
        ":image",
        "//third_party/docker/checkmarx/kics",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
