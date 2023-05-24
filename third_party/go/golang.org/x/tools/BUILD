subinclude("///go//build_defs:go")

go_module(
    name = "tools",
    install = [
        "go/ast/astutil",
        "internal/typeparams",
    ],
    module = "golang.org/x/tools",
    version = "v0.1.10",
    visibility = [
        "//third_party/go/mvdan.cc/gofumpt:all",
    ],
)
