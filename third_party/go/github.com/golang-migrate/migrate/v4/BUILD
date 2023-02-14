go_mod_download(
    name = "_v4_download",
    hashes = ["bb960ab46b4c6b1a6b31bd7b00594c329632b54e7cc6dfe3ba67ef31189f34aa"],
    module = "github.com/golang-migrate/migrate/v4",
    version = "v4.7.0",
)

go_module(
    name = "v4",
    build_tags = ["postgres"],
    download = ":_v4_download",
    install = [
        ".",
        "database",
        "database/postgres",
        "database/stub",
        "internal/...",
        "source",
        "source/file",
        "source/go_bindata",
    ],
    module = "github.com/golang-migrate/migrate/v4",
    visibility = [
        "//pkg/enrichment/db:all",
    ],
    deps = [
        "//third_party/go/github.com/hashicorp/go-multierror",
        "//third_party/go/github.com/lib/pq",
    ],
)

go_module(
    name = "migrate",
    binary = True,
    download = ":_v4_download",
    install = ["cmd/migrate"],
    module = "github.com/golang-migrate/migrate/v4",
    deps = [
        ":v4",
    ],
)
