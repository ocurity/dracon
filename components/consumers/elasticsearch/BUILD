subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "elasticsearch",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "///third_party/go/github.com_elastic_go-elasticsearch_v8//:v8",
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
        "//pkg/templating",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "elasticsearch_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_elastic_go-elasticsearch_v8//:v8",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
        "//pkg/templating",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":elasticsearch"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "elasticsearch",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
