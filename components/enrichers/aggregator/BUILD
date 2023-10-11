subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

kustomized_config(
    name = "k8s",
    srcs = [
        "kustomization.yaml",
        "task.yaml",
    ],
    images = [
        ":image",
    ],
    visibility = ["//examples/pipelines/..."],
)

go_binary(
    name = "aggregator",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "///third_party/go/golang.org_x_crypto//nacl/sign",
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "aggregator_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_google_uuid//:uuid",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/golang.org_x_crypto//nacl/sign",
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":aggregator"],
    visibility = [
        "//examples/...",
    ],
)
