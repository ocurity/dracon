subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

remote_file(
    name = "_dashboard",
    url = "https://github.com/tektoncd/dashboard/releases/download/v0.29.2/tekton-dashboard-release.yaml",
    visibility = ["//build/k8s/tektoncd/dashboard/..."],
)

buildkit_image_mirror(
    name = "image",
    digest = "sha256:67425f5fa32c591578ba6c8b6d4d6a135f8ca195745e023f3f311e3f1ee6ea37",
    repo = "gcr.io/tekton-releases/github.com/tektoncd/dashboard/cmd/dashboard",
    tags = ["v0.29.2"],
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
