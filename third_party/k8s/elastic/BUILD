subinclude("//build/defs:buildkit")

buildkit_image_mirror(
    name = "elasticsearch",
    digest = "sha256:8a84691acd5464776414917e86b6ceb1c8aa38910c20685373aeaab272104619",
    repo = "docker.elastic.co/elasticsearch/elasticsearch",
    tags = [],
    visibility = ["//build/k8s/elastic/..."],
)

buildkit_image_mirror(
    name = "kibana",
    digest = "sha256:34cd9469019f851489aa405211f3f2b3a07092c4aa67fd87e438c45601d6b412",
    repo = "docker.elastic.co/kibana/kibana",
    tags = [],
    visibility = ["//build/k8s/elastic/..."],
)
