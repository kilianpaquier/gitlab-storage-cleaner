# gitlab-storage-cleaner <!-- omit in toc -->

<p align="center">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/kilianpaquier/gitlab-storage-cleaner?include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues-raw/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
  <img alt="GitHub License" src="https://img.shields.io/github/license/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/kilianpaquier/gitlab-storage-cleaner/main?style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/kilianpaquier/gitlab-storage-cleaner/main?style=for-the-badge&label=Go+Version">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
</p>

---

- [How to use ?](#how-to-use-)
  - [Go](#go)
  - [Docker](#docker)
  - [Gitlab CICD](#gitlab-cicd)
  - [Linux](#linux)
- [Commands](#commands)
  - [Artifacts](#artifacts)
  - [Upgrade](#upgrade)

## How to use ?

### Go

```sh
go install github.com/kilianpaquier/gitlab-storage-cleaner/cmd/gitlab-storage-cleaner@latest
```

### Docker

```sh
docker run ghcr.io/kilianpaquier/gitlab-storage-cleaner:v1 artifacts
```

### Gitlab CICD

A potential usage can be to schedule a job once a while with given [`.gitlab-ci.yml`](./.gitlab/.gitlab-ci.yml).

### Linux

```sh
if which gitlab-storage-cleaner >/dev/null; then
  gitlab-storage-cleaner upgrade
  exit $?
fi

OS="linux" # change it depending on our case
ARCH="amd64" # change it depending on our case

echo "installing gitlab-storage-cleaner"
new_version=$(curl -fsSL "https://api.github.com/repos/kilianpaquier/gitlab-storage-cleaner/releases/latest" | jq -r '.tag_name')
url="https://github.com/kilianpaquier/gitlab-storage-cleaner/releases/download/${new_version}/gitlab-storage-cleaner_${OS}_${ARCH}.tar.gz"
curl -fsSL "$url" -o "/tmp/gitlab-storage-cleaner_${OS}_${ARCH}.tar.gz"
mkdir -p "/tmp/gitlab-storage-cleaner/${new_version}"
tar -xzf "/tmp/gitlab-storage-cleaner_${OS}_${ARCH}.tar.gz" -C "/tmp/gitlab-storage-cleaner/${new_version}"
cp "/tmp/gitlab-storage-cleaner/${new_version}/gitlab-storage-cleaner" "${HOME}/.local/bin/gitlab-storage-cleaner"
```

## Commands

```
Usage:
  gitlab-storage-cleaner [command]

Available Commands:
  artifacts   Clean artifacts of provided project(s)' gitlab storage
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  upgrade     Upgrade or install gitlab-storage-cleaner
  version     Show current gitlab-storage-cleaner version

Flags:
  -h, --help                help for gitlab-storage-cleaner
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")

Use "gitlab-storage-cleaner [command] --help" for more information about a command.
```

### Artifacts

```
Clean artifacts of provided project(s)' gitlab storage

Usage:
  gitlab-storage-cleaner artifacts [flags]

Flags:
      --dry-run                       truthy if run must not delete jobs' artifacts but only list matched projects
  -h, --help                          help for artifacts
      --paths strings                 list of valid regexps to match project path (with namespace)
      --server string                 gitlab server host
      --threshold-duration duration   threshold duration (positive) where, from now, jobs artifacts expiration is after will be cleaned up (default 168h0m0s)
      --threshold-size uint           threshold size (in bytes) where jobs artifacts size sum is bigger will be cleaned up (default 1000000)
      --token string                  gitlab read/write token with maintainer rights to delete artifacts

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```

### Upgrade

```
Upgrade or install gitlab-storage-cleaner

Usage:
  gitlab-storage-cleaner upgrade [flags]

Flags:
      --dest string    destination directory where gitlab-storage-cleaner will be upgraded / installed (by default "${HOME}/.local/bin")
  -h, --help           help for upgrade
      --major string   which major version to upgrade / install (must be of the form "v1", "v2", etc.) - mutually exclusive with --minor option
      --minor string   which minor version to upgrade / install (must be of the form "v1.5", "v2.4", etc.) - mutually exclusive with --major option
      --prereleases    whether prereleases are accepted for installation or not

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```