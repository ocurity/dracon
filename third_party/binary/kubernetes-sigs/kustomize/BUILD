VERSION = "v4.5.4"

remote_file(
    name = "kustomize",
    binary = True,
    extract = True,
    hashes = [
        "1159c5c17c964257123b10e7d8864e9fe7f9a580d4124a388e746e4003added3",  # linux-amd64
        "11bfc053dc40c85aa2707f03778b9780011f9ec985251e713e33428f660c8ab4",  # linux_arm64
    ],
    url = f"https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F{VERSION}/kustomize_{VERSION}_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    visibility = ["PUBLIC"],
)
