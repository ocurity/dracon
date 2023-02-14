go_library(
    name = "sarif",
    srcs = [
        "sarif.go",
        "levels.go",
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
        "sarif.go",
        "levels.go",
        "sarif_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/owenrumney/go-sarif/v2",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
    flaky = True,
)
