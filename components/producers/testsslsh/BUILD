subinclude(
    "//build/defs:buildkit",
    "//build/defs:dracon",
)

# this producer covers tfsec https://github.com/aquasecurity/tfsec

dracon_component(
    name = "testsslsh",
    images = [
        ":image",
        "//third_party/docker/drwetter/testssl.sh",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
