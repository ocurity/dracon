subinclude("///go//build_defs:go")

go_library(
    name = "db",
    srcs = [
        "db.go",
        "issue.go",
    ],
    resources = [":resources"],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_golang-migrate_migrate_v4//:v4",
        "///third_party/go/github.com_golang-migrate_migrate_v4//database/postgres",
        "///third_party/go/github.com_golang-migrate_migrate_v4//source/go_bindata",
        "///third_party/go/github.com_jmoiron_sqlx//:sqlx",
        "///third_party/go/github.com_lib_pq//:pq",
        "//api/proto/v1",
        "//third_party/go:google.golang.org_protobuf",
    ],
)

# This is a workaround for embedding files that are
# not in the same pkg dir.
filegroup(
    name = "resources",
    srcs = ["//pkg/enrichment/db/migrations"],
)
