BUILDKIT_VERSION = "0.10.0"

ROOTLESSKIT_VERSION = "1.0.0"

remote_file(
    name = "_buildkit",
    binary = True,
    entry_points = {
        "buildctl": "bin/buildctl",
        "buildkitd": "bin/buildkitd",
    },
    extract = True,
    url = f"https://github.com/moby/buildkit/releases/download/v{BUILDKIT_VERSION}/buildkit-v{BUILDKIT_VERSION}.{CONFIG.OS}-{CONFIG.ARCH}.tar.gz",
)

remote_file(
    name = "_rootlesskit",
    binary = True,
    exported_files = ["rootlesskit"],
    extract = True,
    url = f"https://github.com/rootless-containers/rootlesskit/releases/download/v{ROOTLESSKIT_VERSION}/rootlesskit-x86_64.tar.gz",
)

remote_file(
    name = "_buildctl-daemonless-script",
    binary = True,
    url = f"https://raw.githubusercontent.com/moby/buildkit/v{BUILDKIT_VERSION}/examples/buildctl-daemonless/buildctl-daemonless.sh",
)

subinclude("//build/defs:sh")

sh_script(
    name = "buildctl",
    main = "buildctl_wrapper.sh",
    visibility = ["PUBLIC"],
    deps = [
        ":_buildctl-daemonless-script",
        ":_buildkit",
        ":_rootlesskit",
    ],
)
