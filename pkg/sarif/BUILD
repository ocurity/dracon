subinclude("///go//build_defs:go")

go_library(
    name = "sarif",
    srcs = [
        "levels.go",
        "sarif.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/owenrumney/go-sarif/v2",
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
        "//api/proto/v1",
        "//third_party/go/github.com/owenrumney/go-sarif/v2",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
