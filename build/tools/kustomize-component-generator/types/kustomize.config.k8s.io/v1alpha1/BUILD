subinclude("///go//build_defs:go")

go_library(
    name = "v1alpha1",
    srcs = [
        "component.go",
    ],
    visibility = [
        "//build/tools/kustomize-component-generator/...",
    ],
    deps = [
        "///third_party/go/github.com_mitchellh_go-wordwrap//:go-wordwrap",
        "///third_party/go/gopkg.in_yaml.v3//:yaml.v3",
        "//build/tools/kustomize-component-generator/types/kubernetes",
    ],
)
