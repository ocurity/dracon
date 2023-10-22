subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "docker-trivy",
    images = [
        ":image",
        "//third_party/docker/aquasec/trivy",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
