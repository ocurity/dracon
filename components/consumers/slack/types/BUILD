subinclude("///go//build_defs:go")

go_library(
    name = "types",
    srcs = [
        "full-document.go",
    ],
    visibility = [
        "//components/consumers/slack/...",
    ],
    deps = ["//api/proto/v1"],
)
