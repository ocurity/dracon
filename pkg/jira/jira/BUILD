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
        "///third_party/go/github.com_andygrunwald_go-jira//:go-jira",
        "///third_party/go/github.com_trivago_tgo//tcontainer",
        "//api/proto/v1",
        "//pkg/enumtransformers",
        "//pkg/jira/config",
        "//pkg/jira/document",
        "//pkg/templating",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

go_test(
    name = "jira_test",
    srcs = [
        "api.go",
        "api_test.go",
        "apiutils.go",
        "apiutils_test.go",
    ],
    deps = [
        "///third_party/go/github.com_andygrunwald_go-jira//:go-jira",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_trivago_tgo//tcontainer",
        "//api/proto/v1",
        "//pkg/enumtransformers",
        "//pkg/jira/config",
        "//pkg/jira/document",
        "//pkg/templating",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
