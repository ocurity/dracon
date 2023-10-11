subinclude("///go//build_defs:go")

go_library(
    name = "sarif",
    srcs = [
        "levels.go",
        "sarif.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_owenrumney_go-sarif_v2//sarif",
        "//api/proto/v1",
    ],
)

go_test(
    name = "sarif_test",
    srcs = [
        "levels.go",
        "sarif.go",
        "sarif_test.go",
    ],
    flaky = True,
    deps = [
        "///third_party/go/github.com_owenrumney_go-sarif_v2//sarif",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
