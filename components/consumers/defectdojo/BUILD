subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "defectdojo",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//components/consumers/defectdojo/client",
        "//components/consumers/defectdojo/types",
    ],
)

go_test(
    name = "defectdojo_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/consumers",
        "//components/consumers/defectdojo/client",
        "//components/consumers/defectdojo/types",
        "//third_party/go/github.com/stretchr/testify",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":defectdojo"],
)

dracon_component(
    name = "defectdojo",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
