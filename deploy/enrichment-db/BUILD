subinclude("///k8s//build/defs:helm")

helm_release(
    name = "postgresql",
    kube_context = "//build/k8s/k3d:dracon",
    labels = ["support"],
    namespace = "dracon",
    src_chart = "//third_party/k8s/bitnami/postgresql:postgresql",
    values_files = ["values.yaml"],
)
