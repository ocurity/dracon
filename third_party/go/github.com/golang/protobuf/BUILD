go_mod_download(
    name = "protobuf_download",
    module = "github.com/golang/protobuf",
    version = "v1.5.2",
)

go_module(
    name = "protobuf",
    download = ":protobuf_download",
    install = ["..."],
    module = "github.com/golang/protobuf",
    visibility = [
        # really trying to restrict the usages of this as it is the older API.
        "//third_party/go/google.golang.org/grpc:all",
        "//third_party/go/github.com/envoyproxy/go-control-plane:all",
        "//third_party/go/github.com/envoyproxy/protoc-gen-validate:all",
        "//third_party/go/github.com/cncf/xds/go:all",
        # "//pkg/jira/document:all",
        # "//components/producers/java_spotbugs:all",
        # "//components/consumers/slack/utils:all",
        # "//pkg/enrichment/db:all",
        # "//pkg/enrichment:all",
        # "//components/consumers/slack:all",
        # "//components/consumers/jira_c/utils:all",
        # "//components/consumers/jira_c:all",
    ],
    deps = ["//third_party/go/google.golang.org/protobuf"],
)
