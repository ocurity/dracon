subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "blackduck-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/context",
        "//third_party/go/github.com/package-url/packageurl-go",
    ],
)

go_test(
    name = "blackduck_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/context",
        "//pkg/testutil",
        "//third_party/go/github.com/package-url/packageurl-go",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":blackduck-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "blackduck",
    images = [
        ":image",
        "//third_party/docker/securego/gosec",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
