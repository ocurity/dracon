go_library(
    name = "utils",
    srcs = [
        "utils.go",
    ],
    visibility = [
        "//components/consumers/slack/...",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers/slack/types",
    ],
)

go_test(
    name = "slack_go_test",
    srcs = [
        "utils_test.go",
    ],
    deps = [
        ":utils",
        "//api/proto/v1",
        "//components/consumers",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
