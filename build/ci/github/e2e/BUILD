subinclude("//build/defs:sh")

sh_script(
    name = "wait_for_pipelineruns",
    main = "wait_for_pipelineruns.sh",
    deps = [
        "//build/util",
        "//third_party/binary/mikefarah/yq",
    ],
)
