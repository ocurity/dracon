subinclude("///go//build_defs:go")

go_library(
    name = "v1beta1",
    srcs = [
        "pipeline.go",
        "task.go",
    ],
    visibility = [
        "//build/tools/kustomize-component-generator/...",
    ],
    deps = [
        "///third_party/go/gopkg.in_yaml.v3//:yaml.v3",
        "//build/tools/kustomize-component-generator/types/kubernetes",
    ],
)
