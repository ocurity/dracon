subinclude(
    "///buildkit//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "trufflehog",
    images = [
        ":image",
        "//third_party/docker/trufflesecurity/trufflehog",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
