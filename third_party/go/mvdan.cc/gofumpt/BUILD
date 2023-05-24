subinclude("///go//build_defs:go")

go_module(
    name = "gofumpt",
    binary = True,
    licences = ["BSD-3-Clause"],
    module = "mvdan.cc/gofumpt",
    version = "v0.3.1",
    visibility = [
        "//build/...",
    ],
    deps = [
        "//third_party/go/github.com/google/go-cmp",
        "//third_party/go/golang.org/x/mod",
        "//third_party/go/golang.org/x/sync",
        "//third_party/go/golang.org/x/sys",
        "//third_party/go/golang.org/x/tools",
    ],
)
