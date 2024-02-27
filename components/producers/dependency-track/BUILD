subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "dependency-track-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/context",
    ],
)

go_test(
    name = "dependency-track_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/context",
        "//pkg/testutil",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":dependency-track-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "dependency-track",
    images = [
        ":image",
        "//third_party/docker/securego/gosec",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
