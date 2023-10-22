subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers tfsec https://github.com/aquasecurity/tfsec

dracon_component(
    name = "terraform-tfsec",
    images = [
        ":image",
        "//third_party/docker/aquasec/tfsec",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
