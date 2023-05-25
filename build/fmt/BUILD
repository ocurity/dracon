subinclude("//build/defs:sh")

sh_script(
    name = "please",
    labels = ["fmt"],
    main = "please.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
    ],
)

sh_script(
    name = "go",
    labels = ["fmt"],
    main = "go.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
        "//third_party/go:mvdan.cc_gofumpt",
        "//third_party/go:toolchain",
    ],
)

sh_script(
    name = "proto",
    labels = ["fmt"],
    main = "proto.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
        "//third_party/binary/bufbuild/buf",
    ],
)
