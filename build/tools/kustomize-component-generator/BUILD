go_binary(
    name = "kustomize-component-generator",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//build/tools/kustomize-component-generator/patches",
        "//build/tools/kustomize-component-generator/types/kubernetes",
        "//build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1",
        "//build/tools/kustomize-component-generator/types/tekton.dev/v1beta1",
        "//third_party/go/github.com/mitchellh/go-wordwrap",
        "//third_party/go/gopkg.in/yaml.v3",
    ],
)
