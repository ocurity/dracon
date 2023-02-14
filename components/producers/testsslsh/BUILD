subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers tfsec https://github.com/aquasecurity/tfsec

go_binary(
    name = "testsslsh-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/testsslsh/types",
        "//pkg/sarif",
    ],
)

go_test(
    name = "testsslsh_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/testsslsh/types",
        "//pkg/sarif",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":testsslsh-parser"],
)

dracon_component(
    name = "testsslsh",
    images = [
        ":image",
        "//third_party/docker/drwetter/testssl.sh",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
