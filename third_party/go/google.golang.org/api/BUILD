go_mod_download(
    name = "api_dl",
    module = "google.golang.org/api",
    version = "v0.114.0",
)

go_module(
    name = "api_googleapi",
    download = ":api_dl",
    install = [
        "googleapi",
        "googleapi/transport",
    ],
    module = "google.golang.org/api",
    visibility = [
        "//third_party/go/github.com/googleapis/gax-go/v2",
    ],
)

go_module(
    name = "api",
    download = ":api_dl",
    install = [
        "bigquery/v2",
        "internal",
        "internal/cert",
        "internal/gensupport",
        "internal/impersonate",
        "internal/third_party/uritemplates",
        "iterator",
        "option",
        "option/internaloption",
        "transport",
        "transport/grpc",
        "transport/http",
        "transport/http/internal/propagation",
    ],
    module = "google.golang.org/api",
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/cloud.google.com/go/compute",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/googleapis/enterprise-certificate-proxy",
        "//third_party/go/github.com/googleapis/gax-go/v2",
        "//third_party/go/go.opencensus.io",
        "//third_party/go/golang.org/x/net",
        "//third_party/go/golang.org/x/oauth2",
        "//third_party/go/google.golang.org/grpc",
    ],
)
