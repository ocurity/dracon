VERSION = "v4.6.3"

remote_file(
    name = "yq",
    binary = True,
    extract = True,
    hashes = [
        "0e105edbb0ebc7c4115c610168f1d6b0ff3ceb38043dac92e18fa0698130d69f",  # linux_amd64
        "bd7f08ee1e9f77c14bdf4c20f88afd3e5573401b31d721d3d8b08f9e6164f322",  # linux_arm64
    ],
    url = f"https://github.com/mikefarah/yq/releases/download/{VERSION}/yq_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    visibility = [
        "//api/openapi/tekton/tests/...",
        "//build/...",
    ],
)
