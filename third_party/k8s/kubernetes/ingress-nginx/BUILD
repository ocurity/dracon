subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_remote_chart(
    name = "ingress-nginx",
    images = [
        ":controller",
    ],
    url = "https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.2.5/ingress-nginx-4.2.5.tgz",
    visibility = ["//build/k8s/kubernetes/ingress-nginx/..."],
)

buildkit_image_mirror(
    name = "controller",
    digest = "sha256:54f7fe2c6c5a9db9a0ebf1131797109bb7a4d91f56b9b362bde2abd237dd1974",
    repo = "registry.k8s.io/ingress-nginx/controller",
    tags = ["v1.3.1"],
)
