subinclude("///k8s//build/defs:helm")

helm_release(
    name = "ingress-nginx",
    kube_context = "//build/k8s/k3d:dracon",
    labels = [
        "support",
        "webhook",
    ],
    namespace = "ingress-nginx",
    src_chart = "//third_party/k8s/kubernetes/ingress-nginx:ingress-nginx",
    values_files = ["values.yaml"],
)
