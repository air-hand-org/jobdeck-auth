import os
from pathlib import Path
import subprocess

import pytest

SCRIPT = Path(__file__).with_name('generated-artifacts.sh')


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


def test_allowed_changes_commit_and_request_push(repo):
    (repo / 'go.mod').write_text('module changed\n')
    license = repo / 'third_party/licenses/linux-amd64/example/LICENSE'
    license.parent.mkdir(parents=True)
    license.write_text('license\n')
    result = run(repo)
    assert result.returncode == 0, result.stderr
    assert 'push=true' in result.stdout
    assert git(repo, 'status', '--porcelain') == ''
    assert git(repo, 'log', '-1', '--format=%s') == 'chore: regenerate Renovate artifacts'


@pytest.mark.parametrize('path', ['README.md', 'internal/app/auth_usecase.go', 'third_party/licenses/other/LICENSE'])
def test_outside_change_rejects_entire_commit(repo, path):
    before = git(repo, 'rev-parse', 'HEAD')
    (repo / 'go.mod').write_text('module changed\n')
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
    (repo / 'go.mod').write_text('module changed\n')
    assert run(repo).returncode == 1


@pytest.mark.parametrize('eligible', ['false', ''])
def test_other_prs_fail_without_commit(repo, eligible):
    before = git(repo, 'rev-parse', 'HEAD')
    (repo / 'go.mod').write_text('module changed\n')
    assert run(repo, ELIGIBLE=eligible).returncode == 1
    assert git(repo, 'rev-parse', 'HEAD') == before


def test_wrong_checkout_rejected(repo):
    (repo / 'go.mod').write_text('module changed\n')
    assert run(repo, EXPECTED_HEAD='0' * 40).returncode == 1


@pytest.mark.parametrize('advanced', [False, True])
def test_push_fails_current_run_and_preserves_concurrent_head(repo, tmp_path, advanced):
    remote = repo.parent / (repo.name + '-remote.git')
    git(repo, 'init', '--bare', str(remote))
    git(repo, 'remote', 'add', 'origin', str(remote))
    git(repo, 'push', '-q', 'origin', 'HEAD:refs/heads/renovate/test')
    expected = git(repo, 'rev-parse', 'HEAD')
    (repo / 'go.mod').write_text('module regenerated\n')
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
        env=dict(os.environ, PUSH_TOKEN='fixture', EXPECTED_HEAD=expected, HEAD_BRANCH='renovate/test'),
        capture_output=True, text=True,
    )
    assert result.returncode == 1
    actual = git(repo, 'ls-remote', 'origin', 'refs/heads/renovate/test').split()[0]
    assert actual == (concurrent if advanced else git(repo, 'rev-parse', 'HEAD'))
    assert ('head advanced' if advanced else 'fails intentionally') in result.stdout + result.stderr
