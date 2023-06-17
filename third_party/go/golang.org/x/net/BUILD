subinclude("///go//build_defs:go")

go_module(
    name = "net",
    install = ["..."],
    module = "golang.org/x/net",
    strip = [
        "http2/h2demo",
        "http2/h2i",
    ],
    version = "v0.0.0-20220412020605-290c469a71a5",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/golang.org/x/sys",
        "//third_party/go/golang.org/x/text",
    ],
)
