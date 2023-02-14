subinclude("///k8s//build/defs:k8s")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "my-component-task.yaml",
        "patch-pipeline.yaml",
    ],
    visibility = ["//api/openapi/tekton/tests/pipeline/dag/component-patch/..."],
)
