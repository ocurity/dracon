VERSION = "5.4.4"

remote_file(
    name = "k3d",
    binary = True,
    hashes = [
        "786151745379e511bd4dd95d1593241bb2b0a1d91e088573a1a1cb104732a3bd",  # linux-amd64
        "cf12987491f3092884f71e263f58dbe818dadba40937ae99b1e71a1a83548ef0",  # linux_arm64
    ],
    url = f"https://github.com/rancher/k3d/releases/download/v{VERSION}/k3d-{CONFIG.OS}-{CONFIG.ARCH}",
    visibility = ["//build/k8s/k3d/..."],
)
