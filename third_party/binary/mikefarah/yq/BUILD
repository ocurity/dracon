VERSION = "v4.6.3"

remote_file(
    name = "yq",
    binary = True,
    extract = True,
    hashes = [
        "0e105edbb0ebc7c4115c610168f1d6b0ff3ceb38043dac92e18fa0698130d69f",  # linux_amd64
    ],
    url = f"https://github.com/mikefarah/yq/releases/download/{VERSION}/yq_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    visibility = [
        "//api/openapi/tekton/tests/...",
        "//build/...",
    ],
)
