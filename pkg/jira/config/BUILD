subinclude("///go//build_defs:go")

go_library(
    name = "config",
    srcs = [
        "config.go",
        "types.go",
    ],
    visibility = [
        "PUBLIC",
    ],
)

go_test(
    name = "config_test",
    srcs = [
        "config_test.go",
    ],
    deps = [
        ":config",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
