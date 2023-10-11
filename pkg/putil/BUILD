subinclude("///go//build_defs:go")

go_library(
    name = "putil",
    srcs = [
        "load.go",
        "write.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//api/proto/v1",
        "//third_party/go:google.golang.org_protobuf",
    ],
)
