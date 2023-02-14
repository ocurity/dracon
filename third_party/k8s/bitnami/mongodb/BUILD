subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_remote_chart(
    name = "mongodb",
    images = [
        ":image",
    ],
    url = "https://charts.bitnami.com/bitnami/mongodb-13.3.0.tgz",
    visibility = ["PUBLIC"],
)

buildkit_image_mirror(
    name = "image",
    digest = "sha256:425914ec67ffb9927dc44fc70c546350739cb97cbcc3953d8bda4960488e52f0",
    repo = "docker.io/bitnami/mongodb",
    tags = ["6.0.2-debian-11-r1"],
)
