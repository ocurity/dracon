subinclude("///go//build_defs:go")

go_module(
    name = "go",
    install = [
        "civil",
        "internal",
        "internal/detect",
        "internal/fields",
        "internal/optional",
        "internal/trace",
        "internal/uid",
        "internal/version",
    ],
    module = "cloud.google.com/go",
    version = "v0.110.0",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/googleapis/gax-go/v2",
        "//third_party/go/go.opencensus.io",
        "//third_party/go/golang.org/x/oauth2",
        "//third_party/go/golang.org/x/xerrors",
        "//third_party/go/google.golang.org/api",
        "//third_party/go/google.golang.org/genproto",
        "//third_party/go/google.golang.org/grpc",
    ],
)
