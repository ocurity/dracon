subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

go_binary(
    name = "depsdev",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//pkg/cyclonedx",
        "//pkg/putil",
        "//policies/example.gosec:demo-data",
        "//third_party/go/github.com/CycloneDX/cyclonedx-go",
        "//third_party/go/github.com/package-url/packageurl-go",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "depsdev_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/cyclonedx",
        "//pkg/putil",
        "//third_party/go/github.com/CycloneDX/cyclonedx-go",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/package-url/packageurl-go",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":depsdev"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "depsdev",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
