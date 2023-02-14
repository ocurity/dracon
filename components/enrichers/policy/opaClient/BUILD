subinclude(
    "//build/defs:dracon",
)

go_library(
    name = "opaClient",
    srcs = [
        "opaClient.go",
    ],
    visibility = [
        "//components/enrichers/policy/...",
    ],
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

# dracon_component(
#     name = "policy",
#     images = [
#         ":image",
#     ],
#     task = "task.yaml",
#     visibility = ["//examples/pipelines/..."],
# )
