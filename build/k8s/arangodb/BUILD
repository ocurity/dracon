subinclude("///k8s//build/defs:helm")

helm_release(
    name = "operator",
    kube_context = "//build/k8s/k3d:dracon",
    namespace = "arangodb",
    skip_crds = True,
    src_chart = "//third_party/k8s/arangodb/kube-arangodb",
    values_files = ["arangodb_values.yaml"],
)

helm_chart(
    name = "database",
    src_templates = [
        "server.yaml",
        "ingress.yaml",
    ],
    subcharts = [
        "//third_party/k8s/arangodb/kube-arangodb:crds",
    ],
)

helm_release(
    name = "database",
    kube_context = "//build/k8s/k3d:dracon",
    namespace = "arangodb",
    src_chart = ":database",
)
