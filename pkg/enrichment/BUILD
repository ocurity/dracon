go_library(
    name = "enrichment",
    srcs = [
        "issue.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
        "//pkg/enrichment/db",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

go_test(
    name = "enrichment_test",
    srcs = [
        "issue.go",
        "issue_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/enrichment/db",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
