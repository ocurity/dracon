go_module(
    name = "oauth2",
    install = [
        ".",
        "google",
        "internal",
        "jws",
        "jwt",
    ],
    module = "golang.org/x/oauth2",
    version = "v0.0.0-20220411215720-9780585627b5",
    visibility = [
        "//components/consumers/github",
        "//third_party/go/github.com/google/go-github/...",
        "//third_party/go/google.golang.org/grpc:all",
    ],
    deps = [
        "//third_party/go/cloud.google.com/go/compute",
        "//third_party/go/golang.org/x/net",
    ],
)
