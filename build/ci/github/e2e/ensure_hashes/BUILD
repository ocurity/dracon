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
        "///third_party/go/go.mongodb.org_mongo-driver//bson",
        "///third_party/go/go.mongodb.org_mongo-driver//mongo",
        "///third_party/go/go.mongodb.org_mongo-driver//mongo/options",
        "//api/proto/v1",
    ],
)
