subinclude("///go//build_defs:go")

go_library(
    name = "types",
    srcs = [
        "errors.go",
        "types.go",
    ],
    visibility = ["//components/producers/typescript-npm-audit/..."],
    deps = [
        "//api/proto/v1",
    ],
)
