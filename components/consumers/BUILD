subinclude("///go//build_defs:go")

go_library(
    name = "consumers",
    srcs = [
        "consumer.go",
    ],
    visibility = ["//components/consumers/..."],
    deps = [
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "consumers_go_test",
    srcs = [
        "consumer.go",
        "consumer_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
