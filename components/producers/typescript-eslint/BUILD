subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "typescript-eslint-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-eslint/types",
    ],
)

go_test(
    name = "typescript-eslint_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-eslint/types",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":typescript-eslint-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "typescript-eslint",
    images = [
        ":image",
        "//components/producers/typescript-eslint/eslint-wrapper:eslint",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
