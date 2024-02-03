subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "checkmarx-parser",
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
    name = "checkmarx_test",
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
    srcs = [":checkmarx-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "checkmarx",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
