subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

go_binary(
    name = "ossf-scorecard-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
    ],
)

go_test(
    name = "ossf-scorecard_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":ossf-scorecard-parser"],
    visibility = [
        "//examples/...",
    ],
)

buildkit_image(
    name = "scorecard-dracon",
    dockerfile = "Dockerfile",
)

dracon_component(
    name = "ossf-scorecard",
    images = [
        ":image",
        ":scorecard-dracon",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
