subinclude("///go//build_defs:go")

go_binary(
    name = "kustomize-component-generator",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_mitchellh_go-wordwrap//:go-wordwrap",
        "///third_party/go/gopkg.in_yaml.v3//:yaml.v3",
        "//build/tools/kustomize-component-generator/patches",
        "//build/tools/kustomize-component-generator/types/kubernetes",
        "//build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1",
        "//build/tools/kustomize-component-generator/types/tekton.dev/v1beta1",
    ],
)
