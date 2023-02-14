go_library(
    name = "client",
    srcs = [
        "client.go",
    ],
    visibility = ["//components/consumers/defectdojo/..."],
    deps = [
        "//components/consumers/defectdojo/types",
    ],
)

go_test(
    name = "client_test",
    srcs = [
        "client.go",
        "client_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers/defectdojo/types",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
