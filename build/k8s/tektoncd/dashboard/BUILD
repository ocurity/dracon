subinclude("///k8s//build/defs:helm")

helm_release(
    name = "dashboard",
    kube_context = "//build/k8s/k3d:dracon",
    labels = ["support"],
    namespace = "tekton-pipelines",
    src_chart = "//third_party/k8s/tektoncd/dashboard:dashboard",
    values_files = ["values.yaml"],
)
