go_library(
    name = "types",
    srcs = [
        "yarn-issue.go",
    ],
    visibility = ["//components/producers/typescript-yarn-audit/..."],
    deps = [
        "//api/proto/v1",
        "//components/producers",
    ],
)

go_test(
    name = "types_test",
    srcs = [
        "yarn-issue_test.go",
    ],
    deps = [
        ":types",
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
