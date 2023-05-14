go_library(
    name = "templating",
    srcs = [
        "template_description.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
    ],
)

go_test(
    name = "templating_test",
    srcs = [
        "template_description.go",
        "template_description_test.go",
    ],
    flaky = True,
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
