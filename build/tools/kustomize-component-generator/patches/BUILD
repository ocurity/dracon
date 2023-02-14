go_library(
    name = "patches",
    srcs = glob(["*.go"]),
    visibility = [
        "//build/tools/kustomize-component-generator:all",
    ],
    deps = [
        "//build/tools/kustomize-component-generator/types/kubernetes",
        "//build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1",
        "//build/tools/kustomize-component-generator/types/tekton.dev/v1beta1",
        "//third_party/go/gopkg.in/yaml.v3",
    ],
)
