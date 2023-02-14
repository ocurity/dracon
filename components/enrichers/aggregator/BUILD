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
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/putil",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/golang.org/x/crypto",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "aggregator_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/golang.org/x/crypto",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":aggregator"],
    visibility = [
        "//examples/...",
    ],
)
