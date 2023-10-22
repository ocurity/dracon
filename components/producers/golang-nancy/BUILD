subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

dracon_component(
    name = "golang-nancy",
    images = [
        ":image",
        "//third_party/docker/sonatypecommunity/nancy",
        "//third_party/docker:golang",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
