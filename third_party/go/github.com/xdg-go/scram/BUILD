subinclude("///go//build_defs:go")

go_module(
    name = "scram",
    module = "github.com/xdg-go/scram",
    version = "v1.1.1",
    visibility = [
        "//third_party/go/go.mongodb.org/mongo-driver:all",
    ],
    deps = [
        "//third_party/go/github.com/xdg-go/pbkdf2",
        "//third_party/go/github.com/xdg-go/stringprep",
        "//third_party/go/golang.org/x/crypto",
    ],
)
