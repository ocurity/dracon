subinclude("//build/defs:buildkit")

filegroup(
    name = "npm-audit",
    srcs = [
        "index.js",
        "package.json",
        "package-lock.json",
    ],
)

buildkit_image(
    name = "image",
    srcs = [
        ":npm-audit",
    ],
    dockerfile = "Dockerfile-tool-npm-audit",
)
