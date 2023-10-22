subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "typescript-eslint",
    images = [
        ":image",
        "//components/producers/typescript-eslint/eslint-wrapper:eslint",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
