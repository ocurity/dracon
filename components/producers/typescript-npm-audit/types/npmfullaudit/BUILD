subinclude("///go//build_defs:go")

go_library(
    name = "npmfullaudit",
    srcs = [
        "npmfullaudit.go",
    ],
    visibility = ["//components/producers/typescript-npm-audit/..."],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/typescript-npm-audit/types",
    ],
)

go_test(
    name = "npmfullaudit_test",
    srcs = [
        "npmfullaudit_test.go",
    ],
    deps = [
        ":npmfullaudit",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
    ],
)
