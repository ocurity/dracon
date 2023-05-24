subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers trivy https://github.com/aquasecurity/trivy

go_binary(
    name = "trivy-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/docker-trivy/types",
        "//pkg/cyclonedx",
        "//pkg/sarif",
    ],
)

go_test(
    name = "trivy_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/docker-trivy/types",
        "//pkg/cyclonedx",
        "//pkg/sarif",
        "//third_party/go/github.com/stretchr/testify",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":trivy-parser"],
)

dracon_component(
    name = "docker-trivy",
    images = [
        ":image",
        "//third_party/docker/aquasec/trivy",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
