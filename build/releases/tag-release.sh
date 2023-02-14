#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

git fetch origin main --tags &> /dev/null
commit_sha="$(git rev-parse origin/main)"

highest_tag="$(git tag -l --sort=-v:refname | grep ^v | sed 1q)"

major=$(echo "${highest_tag}" | cut -f1 -d. | tr -dc '0-9')
minor=$(echo "${highest_tag}" | cut -f2 -d. | tr -dc '0-9')
patch=$(echo "${highest_tag}" | cut -f3 -d. | tr -dc '0-9')

util::info "The highest current tag is '${highest_tag}' (major: ${major}, minor: ${minor}, patch: ${patch})"
util::info "What kind of release is this? (For guidance, see: https://semver.org/) [major/minor/patch] " 
read -r release_type
case "${release_type}" in
major)
    major=$((major+1))
    minor=0
    patch=0
    ;;
minor)
    minor=$((minor+1))
    patch=0
    ;;
patch)
    patch=$((patch+1))
    ;;
*)
  util::error "Invalid option: '%s'.\n" "${release_type}"
  exit 1
  ;;
esac

new_tag="v${major}.${minor}.${patch}"

util::prompt "The new release will be '${new_tag}'. Is this OK?"

git tag --annotate "${new_tag}" --message "${new_tag}" "${commit_sha}"
git push origin --tags
