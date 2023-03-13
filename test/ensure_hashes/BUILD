subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "dbComparison",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go/go.mongodb.org/mongo-driver",
    ],
)

filegroup(
    name = "hashes",
    srcs = ["hashes.txt"],
)

buildkit_distroless_image(
    name = "image",
    srcs = [
        ":dbComparison",
    ],
)

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "pod.yaml",
        ":hashes",
    ],
    images = [":image"],
    visibility = ["//examples/pipelines/..."],
)
