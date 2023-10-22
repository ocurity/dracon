subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

dracon_component(
    name = "golang-gosec",
    images = [
        ":image",
        "//third_party/docker/securego/gosec",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
