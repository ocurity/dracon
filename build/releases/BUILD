subinclude("//build/defs:sh")

sh_script(
    name = "pre-release",
    main = "pre-release.sh",
    deps = [
        "//build/util",
    ],
)

sh_script(
    name = "tag-release",
    main = "tag-release.sh",
    deps = [
        "//build/util",
    ],
)

sh_script(
    name = "publish-images",
    main = "publish-images.sh",
    deps = [
        "//build/util",
    ],
)

sh_script(
    name = "remove-images",
    main = "remove-images.sh",
    deps = [
        "//build/util",
    ],
)
