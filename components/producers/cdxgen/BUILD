subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "cdxgen-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/cyclonedx",
        "//pkg/sarif",
    ],
)

go_test(
    name = "cdxgen_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/cyclonedx",
        "//pkg/sarif",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":cdxgen-parser"],
)

dracon_component(
    name = "cdxgen",
    images = [
        ":image",
        "//third_party/docker/cyclonedx/cdxgen",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
