subinclude(
    "//build/defs:buildkit",
    "///k8s//build/defs:helm",
)

VERSION = "2.2.0"

remote_file(
    name = "_crds",
    hashes = ["6e2581aec28871cd483ff08058113262df13df893743e4b865b0fe13d176f8af"],
    url = f"https://download.elastic.co/downloads/eck/{VERSION}/crds.yaml",
    visibility = ["//build/k8s/..."],
)

remote_file(
    name = "_operator",
    hashes = ["6d920700b3b2f9450757baa853e25f52713e981239353de4685d552c78aba9a8"],
    url = f"https://download.elastic.co/downloads/eck/{VERSION}/operator.yaml",
    visibility = ["//build/k8s/..."],
)

buildkit_image_mirror(
    name = "operator",
    digest = "sha256:b633380f6317e6957c90fd144c041726b82810292e0e53505e5df25494dec68e",
    repo = "docker.elastic.co/eck/eck-operator",
    tags = [],
)

helm_chart(
    name = "eck",
    images = [
        ":operator",
    ],
    src_crds = [
        ":_crds",
    ],
    src_templates = [
        ":_operator",
    ],
    visibility = ["//build/k8s/elastic/..."],
)

helm_release(
    name = "operator",
    kube_context = "//build/k8s/k3d:dracon",
    labels = [
        "support",
        "operator",
    ],
    namespace = "eck",
    src_chart = ":eck",
)
