subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_chart(
    name = "dashboard",
    images = [":image"],
    src_templates = [
        ":_dashboard",
        "ingress.yaml",
    ],
    src_values = "values.yaml",
    visibility = ["//build/k8s/tektoncd/dashboard/..."],
)
