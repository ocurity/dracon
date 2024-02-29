subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "enricher",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/enrichment/db",
        "//pkg/putil",
        "//third_party/go/github.com/spf13/cobra",
        "//third_party/go/github.com/spf13/viper",
    ],
)

go_test(
    name = "enricher_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/enrichment/db",
        "//pkg/putil",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/spf13/cobra",
        "//third_party/go/github.com/spf13/viper",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":enricher"],
    visibility = [
        "//examples/...",
    ],
)

# kustomized_config(
#     name = "k8s",
#     srcs = [
#         "kustomization.yaml",
#         "task.yaml",
#     ],
#     images = [
#         ":image",
#     ],
#     visibility = ["//examples/pipelines/..."],
# )

dracon_component(
    name = "multitool-deduplication",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
