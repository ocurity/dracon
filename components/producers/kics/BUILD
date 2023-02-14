subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers tfsec https://github.com/aquasecurity/tfsec

go_binary(
    name = "kics-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/kics/types",
        "//pkg/sarif",
    ],
)

go_test(
    name = "kics_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/kics/types",
        "//pkg/sarif",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":kics-parser"],
)

dracon_component(
    name = "kics",
    images = [
        ":image",
        "//third_party/docker/checkmarx/kics",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
