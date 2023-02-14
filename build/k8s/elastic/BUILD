subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_chart(
    name = "eck",
    src_templates = [
        "elasticsearch.yaml",
        "ingress.yaml",
        "kibana.yaml",
    ],
)

helm_release(
    name = "eck",
    kube_context = "//build/k8s/k3d:dracon",
    labels = ["support"],
    namespace = "dracon",
    src_chart = ":eck",
)
