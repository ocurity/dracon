VERSION = "0.10.0"

remote_file(
    name = "crane",
    binary = True,
    exported_files = ["crane"],
    extract = True,
    url = f"https://github.com/google/go-containerregistry/releases/download/v{VERSION}/go-containerregistry_Linux_x86_64.tar.gz",
    visibility = ["//build/k8s/..."],
)
