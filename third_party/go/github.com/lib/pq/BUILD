subinclude("///go//build_defs:go")

go_module(
    name = "pq",
    hashes = ["8fd2b25ed50efa405fc127081cac8d9c272d02a8d941a22a0c28524095569a2d"],
    install = [
        ".",
        "oid",
        "scram",
    ],
    module = "github.com/lib/pq",
    version = "v1.10.3",
    visibility = [
        "//pkg/enrichment/db:all",
        "//third_party/go/github.com/golang-migrate/migrate/v4:all",
    ],
)
