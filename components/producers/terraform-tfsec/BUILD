subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers tfsec https://github.com/aquasecurity/tfsec

go_binary(
    name = "tfsec-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/terraform-tfsec/types",
        "//pkg/sarif",
    ],
)

go_test(
    name = "tfsec_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/terraform-tfsec/types",
        "//pkg/sarif",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":tfsec-parser"],
)

dracon_component(
    name = "terraform-tfsec",
    images = [
        ":image",
        "//third_party/docker/aquasec/tfsec",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
