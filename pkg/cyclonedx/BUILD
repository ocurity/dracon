subinclude("///go//build_defs:go")

go_library(
    name = "cyclonedx",
    srcs = [
        "cyclonedx.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/CycloneDX/cyclonedx-go",
    ],
)

go_test(
    name = "cyclonedx_test",
    srcs = [
        "cyclonedx.go",
        "cyclonedx_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/CycloneDX/cyclonedx-go",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
