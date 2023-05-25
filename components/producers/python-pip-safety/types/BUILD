subinclude("///go//build_defs:go")

go_library(
    name = "types",
    srcs = [
        "safety-issue.go",
    ],
    visibility = ["//components/producers/python-pip-safety/..."],
)

go_test(
    name = "types-test",
    srcs = [
        "safety-issue.go",
        "safety-issue_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
