VERSION = "v0.8.0"

remote_file(
    name = "shellcheck",
    binary = True,
    exported_files = [f"shellcheck-{VERSION}/shellcheck"],
    extract = True,
    hashes = [
        "ab6ee1b178f014d1b86d1e24da20d1139656c8b0ed34d2867fbb834dad02bf0a",
    ],
    url = f"https://github.com/koalaman/shellcheck/releases/download/{VERSION}/shellcheck-{VERSION}.linux.x86_64.tar.xz",
    visibility = ["PUBLIC"],
)
