subinclude("///go//build_defs:go")

go_library(
    name = "enrichment",
    srcs = [
        "issue.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
        "//pkg/enrichment/db",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "enrichment_test",
    srcs = [
        "issue.go",
        "issue_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//pkg/enrichment/db",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
