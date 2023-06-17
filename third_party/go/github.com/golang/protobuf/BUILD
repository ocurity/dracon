subinclude("///go//build_defs:go")

go_mod_download(
    name = "protobuf_download",
    module = "github.com/golang/protobuf",
    version = "v1.5.2",
)

go_module(
    name = "protobuf",
    download = ":protobuf_download",
    install = ["..."],
    module = "github.com/golang/protobuf",
    visibility = ["PUBLIC"],
    deps = ["//third_party/go/google.golang.org/protobuf"],
)
