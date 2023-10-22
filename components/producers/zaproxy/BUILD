subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

dracon_component(
    name = "zap",
    images = [
        ":image",
        "//third_party/docker/owasp/zap2docker-stable",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
