subinclude("//build/defs:sh")

sh_script(
    name = "please",
    labels = ["lint"],
    main = "please.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
    ],
)

sh_script(
    name = "sh",
    labels = ["lint"],
    main = "sh.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
        "//third_party/binary/koalaman/shellcheck",
    ],
)

sh_script(
    name = "golangci-lint",
    main = "golangci-lint.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util:golangci-lint",
    ],
)

sh_script(
    name = "github_actions",
    labels = ["lint"],
    main = "github_actions.sh",
    visibility = ["PUBLIC"],
    deps = [
        "//build/util",
        "//third_party/binary/rhysd/actionlint",
    ],
)
