subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:k8s",
    "//build/defs:dracon",
)

go_binary(
    name = "dbComparison",
    srcs = [
        "main.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go/go.mongodb.org/mongo-driver",
    ],
)
