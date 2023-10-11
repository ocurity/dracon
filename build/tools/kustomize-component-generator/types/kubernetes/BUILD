subinclude("///go//build_defs:go")

go_library(
    name = "kubernetes",
    srcs = [
        "gvk.go",
        "metadata.go",
    ],
    visibility = [
        "//build/tools/kustomize-component-generator/...",
    ],
)
