subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "bigquery",
    srcs = [
        "main.go",
    ],
    static = True,
    deps = [
        "///third_party/go/cloud.google.com_go_bigquery//:bigquery",
        "///third_party/go/golang.org_x_oauth2//:oauth2",
        "///third_party/go/google.golang.org_api//option",
        "//api/proto/v1",
        "//components/consumers",
        "//pkg/enumtransformers",
    ],
)

buildkit_distroless_image(
    name = "image",
    srcs = [":bigquery"],
    visibility = [
        "//examples/...",
    ],
)

dracon_component(
    name = "bigquery",
    images = [
        ":image",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
