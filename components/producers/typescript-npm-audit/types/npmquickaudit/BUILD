go_library(
    name = "npmquickaudit",
    srcs = [
        "npm_advisory.go",
        "npmquickaudit.go",
    ],
    visibility = ["//components/producers/typescript-npm-audit/..."],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-npm-audit/types",
    ],
)

go_test(
    name = "npmquickaudit_test",
    srcs = [
        "npm_advisory_test.go",
        "npmquickaudit_test.go",
    ],
    data = [
        ":npmquickaudit_test_data",
    ],
    deps = [
        ":npmquickaudit",
        "//api/proto/v1",
        "//components/producers",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/gopkg.in/h2non/gock.v1",
    ],
)

filegroup(
    name = "npmquickaudit_test_data",
    srcs = [
        "npm_advisory_1556",
        "npm_advisory_no_advisorydata",
        "npm_advisory_not_json",
    ],
    visibility = [":npmquickaudit_test"],
)
