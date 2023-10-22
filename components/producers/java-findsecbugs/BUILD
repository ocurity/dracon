subinclude(
    "///buildkit//build/defs:buildkit",
    "//build/defs:dracon",
)

dracon_component(
    name = "java-findsecbugs",
    images = [
        ":image",
        "//third_party/docker/captainfoobar/findsecbugs",
    ],
    task = "task.yaml",
    visibility = ["//examples/pipelines/..."],
)
