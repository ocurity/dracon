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
        "///third_party/go/github.com_mitchellh_mapstructure//:mapstructure",
        "//api/proto/v1",
        "//components/producers",
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
        "///third_party/go/github.com_mitchellh_mapstructure//:mapstructure",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
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
