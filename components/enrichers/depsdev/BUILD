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
        "///third_party/go/github.com_CycloneDX_cyclonedx-go//:cyclonedx-go",
        "///third_party/go/github.com_package-url_packageurl-go//:packageurl-go",
        "//api/proto/v1",
        "//pkg/cyclonedx",
        "//pkg/putil",
        "//policies/example.gosec:demo-data",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "depsdev_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_CycloneDX_cyclonedx-go//:cyclonedx-go",
        "///third_party/go/github.com_google_uuid//:uuid",
        "///third_party/go/github.com_package-url_packageurl-go//:packageurl-go",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//pkg/cyclonedx",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
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
