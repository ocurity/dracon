#!/usr/bin/env bash
set -Exeuo pipefail

reporoot="$(pwd | sed 's#/plz-out/.*##g')"
mkdir -p "${reporoot}/plz-out/buildkit"
temp_dir="$(mktemp -d -p "${reporoot}/plz-out/buildkit" tmp.XXXXX)"
function cleanup {
    chmod -s "${temp_dir}" || true
    chmod -t "${temp_dir}" || true
    chmod 700 "${temp_dir}" || true
    chmod -R 700 "${temp_dir}" || true
    rm -rf "${temp_dir}" || true
}
trap cleanup EXIT

export XDG_RUNTIME_DIR="${temp_dir}/xdg"
mkdir -p "${temp_dir}/xdg"
mkdir -p "${temp_dir}/rootlesskit"
mkdir -p "${temp_dir}/buildkitd"

ROOTLESSKIT=":_rootlesskit"
BUILDKIT=":_buildkit"
BUILDCTL_DAEMONLESS=":_buildctl-daemonless-script"

export BUILDCTL="${reporoot}/$BUILDKIT/buildctl"
export ROOTLESSKIT="${reporoot}/$ROOTLESSKIT --state-dir ${temp_dir}/rootlesskit"
export BUILDKITD="${reporoot}/$BUILDKIT/buildkitd --oci-worker=true --containerd-worker=false --root ${temp_dir}/buildkitd"
"${reporoot}/$BUILDCTL_DAEMONLESS" "$@"
