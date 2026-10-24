import os
from pathlib import Path
import subprocess
import textwrap

import pytest

SCRIPT = Path(__file__).with_name('generated-artifacts.sh')


@pytest.mark.parametrize('event,author,head_repo,eligible', [
    ('pull_request', 'renovate[bot]', 'air-hand-org/jobdeck-auth', True),
    ('pull_request', 'human', 'air-hand-org/jobdeck-auth', False),
    ('pull_request', 'renovate[bot]', 'fork/jobdeck-auth', False),
    ('push', 'renovate[bot]', 'air-hand-org/jobdeck-auth', False),
    ('merge_group', '', '', False),
    ('pull_request_target', 'renovate[bot]', 'air-hand-org/jobdeck-auth', False),
])
def test_workflow_eligibility_and_checkout(tmp_path, event, author, head_repo, eligible):
    workflow = Path(__file__).parents[1] / 'workflows/_test.yaml'
    source = workflow.read_text()
    step = source.split('        id: eligibility\n', 1)[1].split('      - uses: actions/checkout@', 1)[0]
    shell = textwrap.dedent(step.split('        run: |\n', 1)[1])
    output = tmp_path / 'outputs'
    result = subprocess.run(['bash', '-euo', 'pipefail', '-c', shell], env=dict(
        os.environ, EVENT_NAME=event, PR_AUTHOR=author, HEAD_REPO=head_repo,
        GITHUB_REPOSITORY='air-hand-org/jobdeck-auth', HEAD_SHA='actual-head',
        GITHUB_SHA='event-sha', GITHUB_OUTPUT=str(output),
    ), capture_output=True, text=True)
    assert result.returncode == 0, result.stderr
    assert dict(line.split('=', 1) for line in output.read_text().splitlines()) == {
        'eligible': str(eligible).lower(), 'ref': 'actual-head' if eligible else 'event-sha',
    }
    assert 'ref: ${{ steps.eligibility.outputs.ref }}' in source
    assert "if: steps.eligibility.outputs.eligible == 'true'" in source
    assert 'ELIGIBLE: ${{ steps.eligibility.outputs.eligible }}' in source
    assert 'APP_SLUG: ${{ steps.app-token.outputs.app-slug }}' in source


def git(root, *args):
    return subprocess.check_output(['git', '-C', str(root), *args], text=True).strip()


@pytest.fixture
def repo(tmp_path):
    git(tmp_path, 'init', '-q')
    git(tmp_path, 'config', 'user.email', 'test@example.invalid')
    git(tmp_path, 'config', 'user.name', 'Test')
    (tmp_path / 'go.mod').write_text('module example\n')
    git(tmp_path, 'add', 'go.mod')
    git(tmp_path, 'commit', '-qm', 'initial')
    return tmp_path


def run(repo, **overrides):
    env = dict(os.environ, ELIGIBLE='true', EXPECTED_HEAD=git(repo, 'rev-parse', 'HEAD'))
    env.update(overrides)
    return subprocess.run(['bash', str(SCRIPT)], cwd=repo, env=env, capture_output=True, text=True)


def change_license(repo):
    license = repo / 'third_party/licenses/linux-amd64/example/LICENSE'
    license.parent.mkdir(parents=True, exist_ok=True)
    license.write_text('license\n')


def test_allowed_changes_stage_and_request_push(repo):
    before = git(repo, 'rev-parse', 'HEAD')
    change_license(repo)
    result = run(repo)
    assert result.returncode == 0, result.stderr
    assert 'push=true' in result.stdout
    assert git(repo, 'rev-parse', 'HEAD') == before
    assert git(repo, 'diff', '--cached', '--name-only') == 'third_party/licenses/linux-amd64/example/LICENSE'


@pytest.mark.parametrize('path', [
    'README.md', 'internal/app/auth_usecase.go', 'third_party/licenses/other/LICENSE',
    'go.mod', 'go.sum', 'internal/repository/gen/auth.sql.go',
    'internal/repository/gen/db.go', 'internal/repository/gen/models.go',
    'internal/repository/gen/querier.go', 'internal/app/gen_auth_repository_moq_test.go',
    'internal/app/gen_auth_store_moq_test.go', 'internal/app/gen_token_issuer_moq_test.go',
])
def test_outside_change_rejects_entire_commit(repo, path):
    before = git(repo, 'rev-parse', 'HEAD')
    change_license(repo)
    outside = repo / path
    outside.parent.mkdir(parents=True, exist_ok=True)
    outside.write_text('unexpected\n')
    assert run(repo).returncode == 1
    assert git(repo, 'rev-parse', 'HEAD') == before
    assert git(repo, 'diff', '--cached', '--name-only') == ''


def test_no_diff_continues(repo):
    assert run(repo).returncode == 0
    assert 'push=true' not in run(repo).stdout


def test_repeated_generation_stops(repo):
    git(repo, 'commit', '--allow-empty', '-qm', 'chore: regenerate Renovate artifacts')
    change_license(repo)
    assert run(repo).returncode == 1


@pytest.mark.parametrize('eligible', ['false', ''])
def test_other_prs_fail_without_commit(repo, eligible):
    before = git(repo, 'rev-parse', 'HEAD')
    change_license(repo)
    assert run(repo, ELIGIBLE=eligible).returncode == 1
    assert git(repo, 'rev-parse', 'HEAD') == before


def test_wrong_checkout_rejected(repo):
    change_license(repo)
    assert run(repo, EXPECTED_HEAD='0' * 40).returncode == 1


@pytest.mark.parametrize('advanced', [False, True])
def test_push_fails_current_run_and_preserves_concurrent_head(repo, tmp_path, advanced):
    remote = repo.parent / (repo.name + '-remote.git')
    git(repo, 'init', '--bare', str(remote))
    git(repo, 'remote', 'add', 'origin', str(remote))
    git(repo, 'push', '-q', 'origin', 'HEAD:refs/heads/renovate/test')
    expected = git(repo, 'rev-parse', 'HEAD')
    change_license(repo)
    assert run(repo).returncode == 0
    if advanced:
        other = tmp_path / 'other'
        git(repo, 'clone', '-q', '-b', 'renovate/test', str(remote), str(other))
        git(other, 'config', 'user.email', 'other@example.invalid')
        git(other, 'config', 'user.name', 'Other')
        git(other, 'commit', '--allow-empty', '-qm', 'concurrent change')
        git(other, 'push', '-q', 'origin', 'HEAD')
        concurrent = git(other, 'rev-parse', 'HEAD')
    result = subprocess.run(
        ['bash', str(SCRIPT.with_name('push-generated-artifacts.sh'))], cwd=repo,
        env=dict(os.environ, PUSH_TOKEN='fixture', APP_SLUG='configured-release-app',
                 EXPECTED_HEAD=expected, HEAD_BRANCH='renovate/test'),
        capture_output=True, text=True,
    )
    assert result.returncode == 1
    actual = git(repo, 'ls-remote', 'origin', 'refs/heads/renovate/test').split()[0]
    assert actual == (concurrent if advanced else git(repo, 'rev-parse', 'HEAD'))
    assert ('head advanced' if advanced else 'fails intentionally') in result.stdout + result.stderr
    if not advanced:
        assert git(repo, 'log', '-1', '--format=%an <%ae>|%cn <%ce>') == (
            'configured-release-app[bot] <configured-release-app[bot]@users.noreply.github.com>|'
            'configured-release-app[bot] <configured-release-app[bot]@users.noreply.github.com>'
        )
        assert git(repo, 'rev-parse', 'HEAD^') == expected
        assert git(repo, 'log', '-1', '--format=%s') == 'chore: regenerate Renovate artifacts'
