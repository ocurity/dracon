subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_remote_chart(
    name = "kube-arangodb",
    url = "https://github.com/arangodb/kube-arangodb/releases/download/1.2.19/kube-arangodb-1.2.19.tgz",
    visibility = ["//build/k8s/arangodb/..."],
)

helm_remote_chart(
    name = "crds",
    url = "https://github.com/arangodb/kube-arangodb/releases/download/1.2.19/kube-arangodb-crd-1.2.19.tgz",
    visibility = ["//build/k8s/arangodb/..."],
)
