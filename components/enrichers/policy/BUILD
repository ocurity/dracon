subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

go_binary(
    name = "policy",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/enrichers/policy/opaClient",
        "//pkg/putil",
        "//policies/example.gosec:demo-data",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "policy_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/enrichers/policy/opaClient",
        "//pkg/putil",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":policy"],
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
    name = "policy",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
