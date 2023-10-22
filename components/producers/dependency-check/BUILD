subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

dracon_component(
    name = "dependency-check",
    images = [
        ":image",
        "//third_party/docker/owasp/dependency-check",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
