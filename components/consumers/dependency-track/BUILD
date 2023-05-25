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
        "///third_party/go/github.com_DependencyTrack_client-go//:client-go",
        "///third_party/go/github.com_google_uuid//:uuid",
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/cyclonedx",
    ],
)

go_test(
    name = "dependency-track_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_DependencyTrack_client-go//:client-go",
        "///third_party/go/github.com_google_uuid//:uuid",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/cyclonedx",
        "//third_party/go:google.golang.org_protobuf",
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
