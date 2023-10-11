subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

# this producer covers nancy https://github.com/sonatype-nexus-community/nancy

go_binary(
    name = "golang-nancy-parser",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/golang-nancy/types",
    ],
)

go_test(
    name = "nancy_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        "//api/proto/v1",
        "//components/producers",
        "//components/producers/golang-nancy/types",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":golang-nancy-parser"],
)

dracon_component(
    name = "golang-nancy",
    images = [
        ":image",
        "//third_party/docker/sonatypecommunity/nancy",
        "//third_party/docker:golang",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
