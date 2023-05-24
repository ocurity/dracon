subinclude("///go//build_defs:go")

go_module(
    name = "go-control-plane",
    install = [
        "envoy/config/core/v3",
        "envoy/config/rbac/v3",
        "envoy/config/route/v3",
        "envoy/type/matcher/v3",
    ],
    module = "github.com/envoyproxy/go-control-plane",
    version = "v0.10.1",
    visibility = [
        "//third_party/go/google.golang.org/grpc:all",
    ],
    deps = [
        "//third_party/go/github.com/cncf/xds/go",
        "//third_party/go/github.com/envoyproxy/protoc-gen-validate",
        "//third_party/go/github.com/golang/protobuf",
        "//third_party/go/google.golang.org/genproto",
    ],
)
