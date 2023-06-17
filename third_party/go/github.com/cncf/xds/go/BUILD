subinclude("///go//build_defs:go")

go_module(
    name = "go",
    install = [
        "udpa/annotations",
        "xds/annotations/v3",
        "xds/core/v3",
        "xds/type/matcher/v3",
    ],
    module = "github.com/cncf/xds/go",
    version = "v0.0.0-20220330162227-eded343319d0",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/envoyproxy/protoc-gen-validate",
        "//third_party/go/github.com/golang/protobuf",
    ],
)
