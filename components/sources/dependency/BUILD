subinclude("///go//build_defs:go")

subinclude(
    "///k8s//build/defs:k8s",
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "dependency",
    srcs = [
        "main.go",
    ],
    deps = [
        "///third_party/go/github.com_package-url_packageurl-go//:packageurl-go",
    ],
)

go_test(
    name = "dependency_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_package-url_packageurl-go//:packageurl-go",
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":dependency"],
)

dracon_component(
    name = "source-dependency",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
