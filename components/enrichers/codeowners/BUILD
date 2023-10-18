subinclude(
    "//build/defs:dracon",
    "//build/defs:buildkit",
)

go_binary(
    name = "codeowners",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//pkg/putil",
        "//third_party/go/github.com/hairyhenderson/go-codeowners",
        "//third_party/go/github.com/package-url/packageurl-go",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":codeowners"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "codeowners",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
