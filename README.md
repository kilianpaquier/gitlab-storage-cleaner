# gitlab-storage-cleaner <!-- omit in toc -->

<div align="center">
  <img alt="GitLab Release" src="https://img.shields.io/gitlab/v/release/kilianpaquier%2Fgitlab-storage-cleaner?gitlab_url=https%3A%2F%2Fgitlab.com&include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitLab Issues" src="https://img.shields.io/gitlab/issues/open/kilianpaquier%2Fgitlab-storage-cleaner?gitlab_url=https%3A%2F%2Fgitlab.com&style=for-the-badge">
  <img alt="GitLab License" src="https://img.shields.io/gitlab/license/kilianpaquier%2Fgitlab-storage-cleaner?gitlab_url=https%3A%2F%2Fgitlab.com&style=for-the-badge">
  <img alt="GitLab CICD" src="https://img.shields.io/gitlab/pipeline-status/kilianpaquier%2Fgitlab-storage-cleaner?gitlab_url=https%3A%2F%2Fgitlab.com&branch=main&style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/gitlab/go-mod/go-version/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/gitlab.com/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
</div>

---

- [How to use ?](#how-to-use-)
  - [Go](#go)
  - [Mise](#mise)
  - [Gitlab CICD](#gitlab-cicd)
  - [Linux](#linux)
- [Commands](#commands)
  - [Artifacts](#artifacts)

## How to use ?

### Go

```sh
go install github.com/kilianpaquier/gitlab-storage-cleaner/cmd/gitlab-storage-cleaner@latest
```

### Mise

```sh
mise use -g github:kilianpaquier/gitlab-storage-cleaner
```

### Gitlab CICD

A potential usage can be to schedule a job once a while with given [`.gitlab-ci.yml`](./.gitlab/.gitlab-ci.yml).

### Linux

```sh
OS="linux" # change it depending on your case
ARCH="amd64" # change it depending on your case
INSTALL_DIR="$HOME/.local/bin" # change it depending on your case

new_version=$(curl -fsSL "https://api.github.com/repos/kilianpaquier/gitlab-storage-cleaner/releases/latest" | jq -r '.tag_name')
url="https://github.com/kilianpaquier/gitlab-storage-cleaner/releases/download/$new_version/gitlab-storage-cleaner_${OS}_${ARCH}.tar.gz"
curl -fsSL "$url" | (mkdir -p "/tmp/gitlab-storage-cleaner/$new_version" && cd "/tmp/gitlab-storage-cleaner/$new_version" && tar -xz)
cp "/tmp/gitlab-storage-cleaner/$new_version/gitlab-storage-cleaner" "$INSTALL_DIR/gitlab-storage-cleaner"
```

## Commands

```
Usage:
  gitlab-storage-cleaner [command]

Available Commands:
  artifacts   Clean artifacts of provided project(s)' gitlab storage
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Show current version

Flags:
  -h, --help                help for gitlab-storage-cleaner
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")

Use "gitlab-storage-cleaner [command] --help" for more information about a command.
```

### Artifacts

```
Usage:
  gitlab-storage-cleaner artifacts [flags]

Flags:
      --dry-run                       truthy if run must not delete jobs' artifacts but only list matched projects
  -h, --help                          help for artifacts
      --paths strings                 list of valid regexps to match project path (with namespace)
      --server string                 gitlab server host
      --threshold-duration duration   threshold duration (positive) where, jobs older than command execution time minus this threshold will be deleted (default 168h0m0s)
      --token string                  gitlab read/write token with maintainer rights to delete artifacts

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```

#### Flags

Both CLI flags can be used and environment variables, while the priority is still given to the CLI flags.

| CLI flag               | Environment variable(s)           | Required |
| ---------------------- | --------------------------------- | -------- |
| `--log-format`         | `LOG_FORMAT`                      | No       |
| `--log-level`          | `LOG_LEVEL`                       | No       |
| `--token`              | `GITLAB_TOKEN`, `GL_TOKEN`        | Yes      |
| `--server`             | `CI_API_V4_URL`, `CI_SERVER_HOST` | Yes      |
| `--dry-run`            | `CLEANER_DRY_RUN`                 | No       |
| `--paths`              | `CLEANER_PATHS`                   | Yes      |
| `--threshold-duration` | `CLEANER_THRESHOLD_DURATION`      | No       |
