subinclude("///go//build_defs:go")

go_library(
    name = "producers",
    srcs = [
        "producer.go",
    ],
    visibility = ["//components/producers/..."],
    deps = [
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "producers_test",
    srcs = [
        "producer.go",
        "producer_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go/github.com/mitchellh/mapstructure",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
