go_library(
    name = "context",
    srcs = [
        "context.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
    ],
)

go_test(
    name = "context_test",
    srcs = [
        "context.go",
        "context_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//pkg/testutil",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
