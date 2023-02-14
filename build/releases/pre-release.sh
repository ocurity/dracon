#!/usr/bin/env bash
set -Eeuo pipefail

source "//build/util"

git fetch --tags --force
version="$(git describe --always)"

github_api_version="application/vnd.github.v3+json"
github_owner="ocurity"
github_repo="dracon"

post_release_json=$(cat <<EOF
{
  "tag_name": "${version}",
  "target_commitish": "main",
  "name": "${version}",
  "body": "",
  "draft": false,
  "prerelease": true
}
EOF
)

release_resp=$(curl \
  --header "Accept: ${github_api_version}" \
  --header "Authorization: token ${GITHUB_TOKEN}" \
  --silent \
  --request POST \
  --data "${post_release_json}" \
  "https://api.github.com/repos/${github_owner}/${github_repo}/releases")

release_id=$(echo "${release_resp}" | jq '.id')
if [ -z "${release_id}" ] || [ "${release_id}" == "null" ]; then
  util::error "Could not find release id in response"
  echo "${release_resp}"
  exit 1
fi
util::success "Created pre-release ${release_id}"
