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
        "//build/tools/kustomize-component-generator/types/kubernetes",
        "//third_party/go/github.com/mitchellh/go-wordwrap",
        "//third_party/go/gopkg.in/yaml.v3",
    ],
)
