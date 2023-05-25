subinclude("///go//build_defs:go")

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
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
    ],
)
