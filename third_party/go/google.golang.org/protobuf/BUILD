go_mod_download(
    name = "protobuf_download",
    module = "google.golang.org/protobuf",
    version = "v1.28.0",
)

go_module(
    name = "protobuf",
    download = ":protobuf_download",
    install = ["..."],
    module = "google.golang.org/protobuf",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/google/go-cmp",
    ],
)

go_module(
    name = "protoc-gen-go",
    binary = True,
    download = ":protobuf_download",
    install = ["cmd/protoc-gen-go"],
    module = "github.com/golang/protobuf",
    visibility = ["PUBLIC"],
    deps = [":protobuf"],
)
