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
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "producers_test",
    srcs = [
        "producer.go",
        "producer_test.go",
    ],
    deps = [
        "///third_party/go/github.com_mitchellh_mapstructure//:mapstructure",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
