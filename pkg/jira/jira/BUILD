subinclude("///go//build_defs:go")

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
        "//api/proto/v1",
        "//pkg/enumtransformers",
        "//pkg/jira/config",
        "//pkg/jira/document",
        "//pkg/templating",
        "//third_party/go/github.com/andygrunwald/go-jira",
        "//third_party/go/github.com/trivago/tgo",
        "//third_party/go/google.golang.org/protobuf",
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
        "//api/proto/v1",
        "//pkg/enumtransformers",
        "//pkg/templating",
        "//third_party/go/github.com/andygrunwald/go-jira",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
