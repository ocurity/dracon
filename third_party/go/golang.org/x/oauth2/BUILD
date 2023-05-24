subinclude("///go//build_defs:go")

go_module(
    name = "oauth2",
    install = [
        ".",
        "authhandler",
        "google",
        "google/internal/externalaccount",
        "internal",
        "jws",
        "jwt",
    ],
    module = "golang.org/x/oauth2",
    version = "v0.6.0",
    visibility = ["PUBLIC"],
    deps = ["//third_party/go/cloud.google.com/go/compute"],
)
