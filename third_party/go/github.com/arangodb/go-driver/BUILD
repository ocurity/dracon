subinclude("///go//build_defs:go")

go_module(
    name = "go-driver",
    install = [
        ".",
        "cluster",
        "http",
        "util",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/arangodb/go-driver",
    version = "v1.4.0",
    visibility = ["//components/consumers/arangodb/..."],
    deps = [
        "//third_party/go/github.com/arangodb/go-velocypack",
        "//third_party/go/github.com/pkg/errors",
    ],
)
