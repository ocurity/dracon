subinclude("///go//build_defs:go")

subinclude(
    "///k8s//build/defs:k8s",
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
    name = "tagger",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "///third_party/go/github.com_google_uuid//:uuid",
        "//api/proto/v1",
        "//pkg/enrichment",
        "//pkg/putil",
    ],
)

go_test(
    name = "tagger_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_google_uuid//:uuid",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":tagger"],
    visibility = [
        "//examples/...",
    ],
)
