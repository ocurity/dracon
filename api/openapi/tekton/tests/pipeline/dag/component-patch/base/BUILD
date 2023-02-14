subinclude("///k8s//build/defs:k8s")

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "pipeline.yaml",
        ":api",
    ],
    visibility = ["//api/openapi/tekton/tests/pipeline/dag/component-patch/..."],
)

filegroup(
    name = "api",
    srcs = [
        "//components/base:openapi_schema",
        "//components/base:transformer",
    ],
)
