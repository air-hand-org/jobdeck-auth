#!/usr/bin/env bash
set -euo pipefail

git check-ref-format "refs/heads/${HEAD_BRANCH:?}"
[[ $(git rev-parse HEAD) == "${EXPECTED_HEAD:?}" ]] || exit 1
# This token is supplied only to these Git operations, never persisted in config.
auth=$(printf 'x-access-token:%s' "${PUSH_TOKEN:?}" | base64 -w0)
remote_head=$(git -c "http.https://github.com/.extraheader=AUTHORIZATION: basic $auth" \
  ls-remote --exit-code origin "refs/heads/$HEAD_BRANCH")
if [[ ${remote_head%%$'\t'*} != "$EXPECTED_HEAD" ]]; then
  echo 'PR head advanced; refusing to overwrite concurrent changes.' >&2
  exit 1
fi
# Use the configured App identity returned by the token action, not a guessed slug.
git -c user.name="${APP_SLUG:?}[bot]" -c user.email="${APP_SLUG}[bot]@users.noreply.github.com" \
  commit -m 'chore: regenerate Renovate artifacts'
# A concurrent update after the check is rejected by normal non-force push.
git -c "http.https://github.com/.extraheader=AUTHORIZATION: basic $auth" \
  push origin "HEAD:refs/heads/$HEAD_BRANCH"
echo 'Artifacts pushed. This run fails intentionally; the new head must pass normal CI.'
exit 1
