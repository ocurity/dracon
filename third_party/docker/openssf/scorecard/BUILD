subinclude("//build/defs:buildkit")

VERSION = "4.7.03"

buildkit_image_mirror(
    name = "scorecard",
    digest = "sha256:f186e26685f2d9adf5a4269331319e257127fb33e96f06b8363b877a20e3651e",
    repo = "gcr.io/openssf/scorecard",
    tags = [VERSION],
    visibility = ["//components/producers/ossf-scorecard/..."],
)
