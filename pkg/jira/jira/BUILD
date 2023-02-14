go_library(
    name = "jira",
    srcs = [
        "api.go",
        "apiutils.go",
    ],
    visibility = [
        "PUBLIC",
    ],
    deps = [
        "//pkg/jira/config",
        "//pkg/jira/document",
        "//third_party/go/github.com/andygrunwald/go-jira",
        "//third_party/go/github.com/trivago/tgo",
    ],
)

go_test(
    name = "jira_test",
    srcs = [
        "api_test.go",
        "apiutils_test.go",
    ],
    deps = [
        ":jira",
        "//third_party/go/github.com/andygrunwald/go-jira",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
