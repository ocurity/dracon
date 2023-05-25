subinclude("///go//build_defs:go")

go_library(
    name = "cyclonedx",
    srcs = [
        "cyclonedx.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_CycloneDX_cyclonedx-go//:cyclonedx-go",
        "//api/proto/v1",
    ],
)

go_test(
    name = "cyclonedx_test",
    srcs = [
        "cyclonedx.go",
        "cyclonedx_test.go",
    ],
    deps = [
        "///third_party/go/github.com_CycloneDX_cyclonedx-go//:cyclonedx-go",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
    ],
)
