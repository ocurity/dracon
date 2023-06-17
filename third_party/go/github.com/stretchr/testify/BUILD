subinclude("///go//build_defs:go")
go_module(
    name = "testify",
    install = [
        "assert",
        "require",
    ],
    module = "github.com/stretchr/testify",
    version = "v1.7.1",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/davecgh/go-spew",
        "//third_party/go/github.com/pmezard/go-difflib",
        "//third_party/go/gopkg.in/yaml.v3",
    ],
)
