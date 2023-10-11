subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "enricher",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "///third_party/go/github.com_spf13_cobra//:cobra",
        "///third_party/go/github.com_spf13_viper//:viper",
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/enrichment/db",
        "//pkg/putil",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":enricher"],
    visibility = [
        "//examples/...",
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

dracon_component(
    name = "deduplication",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
