subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "zaproxy-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/zaproxy/types",
    ],
)

go_test(
    name = "zaproxy_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/zaproxy/types",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":zaproxy-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "zap",
    images = [
        ":image",
        "//third_party/docker/owasp/zap2docker-stable",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
