subinclude("///go//build_defs:go")

subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

go_binary(
    name = "eslint-wrapper",
    srcs = [
        "main.go",
    ],
)

filegroup(
    name = "eslint-config",
    srcs = [
        "eslintrc.js",
        "package.json",
    ],
)

buildkit_image(
    name = "eslint",
    srcs = [
        ":eslint-config",
        ":eslint-wrapper",
    ],
    dockerfile = "eslint.Dockerfile",
    visibility = [
        "//examples/...",
    ],
)
