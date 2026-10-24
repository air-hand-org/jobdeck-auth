#!/usr/bin/env bash
set -euo pipefail

# NUL-delimited names include staged, unstaged, deleted and untracked files.
mapfile -d '' -t paths < <({ git diff --no-renames --name-only -z HEAD; git ls-files --others --exclude-standard -z; } | sort -zu)
if ((${#paths[@]} == 0)); then
  exit 0
fi
git status --short
if [[ ${ELIGIBLE:-false} != true ]]; then
  echo 'Only same-repository Renovate pull requests may regenerate artifacts.' >&2
  exit 1
fi
[[ $(git rev-parse HEAD) == "${EXPECTED_HEAD:?}" ]] || exit 1
if [[ $(git log -1 --format=%s) == 'chore: regenerate Renovate artifacts' ]]; then
  echo 'Generated artifacts changed again; refusing an automatic commit loop.' >&2
  exit 1
fi
for path in "${paths[@]}"; do
  case "$path" in
    go.mod|go.sum|internal/repository/gen/auth.sql.go|internal/repository/gen/db.go|internal/repository/gen/models.go|internal/repository/gen/querier.go|internal/app/gen_auth_repository_moq_test.go|internal/app/gen_auth_store_moq_test.go|internal/app/gen_token_issuer_moq_test.go|third_party/licenses/linux-amd64/*) ;;
    *) printf 'Unexpected generated change: %q\n' "$path" >&2; exit 1 ;;
  esac
done
# Validate every path before staging anything; never stage the entire tree.
git add -- "${paths[@]}"
git -c user.name='jobdeck automation' -c user.email='jobdeck-automation@users.noreply.github.com' \
  commit -m 'chore: regenerate Renovate artifacts'
printf 'push=true\n'
if [[ -n ${GITHUB_OUTPUT:-} ]]; then
  printf 'push=true\n' >> "$GITHUB_OUTPUT"
fi
