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
        "//api/proto/v1",
        "//components/consumers",
        "//third_party/go/github.com/elastic/go-elasticsearch/v8",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "elasticsearch_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//third_party/go/github.com/elastic/go-elasticsearch/v8",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
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
