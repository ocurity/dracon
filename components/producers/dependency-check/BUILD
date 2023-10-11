subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "dependency-check-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
    ],
)

go_test(
    name = "dependency_check_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":dependency-check-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "dependency-check",
    images = [
        ":image",
        "//third_party/docker/owasp/dependency-check",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
