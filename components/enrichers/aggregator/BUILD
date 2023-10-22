subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "task.yaml",
    ],
    images = [
        ":image",
    ],
    visibility = ["//examples/pipelines/..."],
)
