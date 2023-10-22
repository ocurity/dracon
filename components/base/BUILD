subinclude("///k8s//build/defs:k8s")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "pipeline.yaml",
        "task.yaml",
        ":openapi_schema",
        ":transformer",
    ],
    images = [
        "//third_party/docker:busybox",
    ],
    visibility = ["PUBLIC"],
)
