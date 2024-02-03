subinclude(
    "///k8s//build/defs:k8s",
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "source-fetch_scan_results",
    images = [
        "//third_party/docker/library/buildpack-dependencies/stable-curl",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
