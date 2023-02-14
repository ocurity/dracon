subinclude(
    "///k8s//build/defs:k8s",
)

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "task.yaml",
    ],
    images = [
        "//third_party/k8s/tektoncd/catalog:git-init",
    ],
    visibility = ["//examples/pipelines/..."],
)
