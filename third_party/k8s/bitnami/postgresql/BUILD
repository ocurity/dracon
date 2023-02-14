subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

helm_remote_chart(
    name = "postgresql",
    images = [
        ":image",
    ],
    url = "https://charts.bitnami.com/bitnami/postgresql-11.9.8.tgz",
    visibility = ["//build/k8s/enrichment-db/..."],
)

buildkit_image_mirror(
    name = "image",
    digest = "sha256:9a868adde823ba8d9d056f3524825f26b021619a63650bba0a06620a12b7e709",
    repo = "docker.io/bitnami/postgresql",
    tags = ["14.5.0-debian-11-r24"],
)
