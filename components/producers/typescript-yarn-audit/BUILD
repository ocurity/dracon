subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "typescript-yarn-audit",
    images = [
        ":image",
        "//third_party/docker:node",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
