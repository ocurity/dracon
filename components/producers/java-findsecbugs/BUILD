subinclude("///go//build_defs:go")

subinclude(
    "///buildkit//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "java-findsecbugs-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//pkg/sarif",
    ],
)

go_test(
    name = "java-findsecbugs-parser_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
        "//pkg/sarif",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":java-findsecbugs-parser"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "java-findsecbugs",
    images = [
        ":image",
        "//third_party/docker/captainfoobar/findsecbugs",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
