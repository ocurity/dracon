subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "dependency-track",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/cyclonedx",
        "//third_party/go/github.com/DependencyTrack/client-go",
        "//third_party/go/github.com/google/uuid",
    ],
)

go_test(
    name = "dependency-track_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/cyclonedx",
        "//third_party/go/github.com/DependencyTrack/client-go",
        "//third_party/go/github.com/google/uuid",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":dependency-track"],
)

dracon_component(
    name = "dependency-track",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
