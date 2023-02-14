#!/usr/bin/env bash
# This script removes all packages from GHCR.
set -Eeuo pipefail

source "//build/util"

ORGANIZATION="ocurity"

mapfile -t all_org_packages < \
    <(
        curl -SsL \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $GITHUB_TOKEN"\
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "https://api.github.com/orgs/$ORGANIZATION/packages?package_type=container" \
        | jq -r '.[].name'
    )

for org_package in "${all_org_packages[@]}"; do
    util::info "deleting $org_package"
    pkg_name_encoded="$(printf '%s' "$org_package" | jq -sRr @uri)"
    curl \
      -X DELETE \
      -H "Accept: application/vnd.github+json" \
      -H "Authorization: Bearer $GITHUB_TOKEN"\
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "https://api.github.com/orgs/$ORGANIZATION/packages/container/$pkg_name_encoded"
done
