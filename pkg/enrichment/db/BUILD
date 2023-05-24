subinclude("///go//build_defs:go")

go_library(
    name = "db",
    srcs = [
        "db.go",
        "issue.go",
    ],
    resources = [":resources"],
    visibility = [
        "//components/enrichers/...",
        "//components/producers/...",
        "//pkg/enrichment/...",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go/github.com/golang-migrate/migrate/v4",
        "//third_party/go/github.com/jmoiron/sqlx",
        "//third_party/go/github.com/lib/pq",
        "//third_party/go/google.golang.org/protobuf",
    ],
)

# This is a workaround for embedding files that are
# not in the same pkg dir.
filegroup(
    name = "resources",
    srcs = ["//pkg/enrichment/db/migrations"],
)
