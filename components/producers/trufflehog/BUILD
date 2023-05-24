subinclude("///go//build_defs:go")

subinclude(
    "///buildkit//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "trufflehog-parser",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/mitchellh/mapstructure",
    ],
)

go_test(
    name = "trufflehog_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/mitchellh/mapstructure",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":trufflehog-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "trufflehog",
    images = [
        ":image",
        "//third_party/docker/trufflesecurity/trufflehog",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
