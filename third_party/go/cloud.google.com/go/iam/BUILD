subinclude("///go//build_defs:go")
go_module(
    name = "iam",
    install = [
        ".",
        "apiv1/iampb",
    ],
    module = "cloud.google.com/go/iam",
    version = "v0.13.0",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/googleapis/gax-go/v2",
        "//third_party/go/google.golang.org/genproto",
        "//third_party/go/google.golang.org/grpc",
        "//third_party/go/google.golang.org/protobuf",
    ],
)
