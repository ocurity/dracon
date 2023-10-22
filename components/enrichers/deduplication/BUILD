subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

# kustomized_config(
#     name = "k8s",
#     srcs = [
#         "kustomization.yaml",
#         "task.yaml",
#     ],
#     images = [
#         ":image",
#     ],
#     visibility = ["//examples/pipelines/..."],
# )

dracon_component(
    name = "deduplication",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
