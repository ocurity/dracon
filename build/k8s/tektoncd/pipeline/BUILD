subinclude("///k8s//build/defs:helm")

helm_release(
    name = "pipeline",
    kube_context = "//build/k8s/k3d:dracon",
    labels = [
        "support",
        "webhook",
    ],
    namespace = "tekton-pipelines",
    src_chart = "//third_party/k8s/tektoncd/pipeline:pipeline",
)
