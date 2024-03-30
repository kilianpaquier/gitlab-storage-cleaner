<!-- This file is safe to edit. Once it exists it will not be overwritten. -->

# gitlab-storage-cleaner <!-- omit in toc -->

<p align="center">
  <img alt="GitHub Actions" src="https://img.shields.io/github/actions/workflow/status/kilianpaquier/gitlab-storage-cleaner/integration.yml?branch=main&style=for-the-badge">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/kilianpaquier/gitlab-storage-cleaner?include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues-raw/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
  <img alt="GitHub License" src="https://img.shields.io/github/license/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/kilianpaquier/gitlab-storage-cleaner/main?style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/kilianpaquier/gitlab-storage-cleaner/main?style=for-the-badge&label=Go+Version">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kilianpaquier/gitlab-storage-cleaner?style=for-the-badge">
</p>

---

- [How to use ?](#how-to-use-)
  - [Gitlab CICD](#gitlab-cicd)
- [Commands](#commands)
  - [Artifacts](#artifacts)

## How to use ?

```sh
go install github.com/kilianpaquier/gitlab-storage-cleaner/cmd/gitlab-storage-cleaner@latest
```

### Gitlab CICD

A potential usage can be to schedule a job once a while with given [`.gitlab-ci.yml`](./.gitlab/.gitlab-ci.yml).

## Commands

```
gitlab-storage-cleaner stands here to help in a small gitlab 
continuous integration's step to clean all old or outdated storage of a given 
or even multiple projects.

Usage:
  gitlab-storage-cleaner [command]

Available Commands:
  artifacts   Clean artifacts of provided project(s)' gitlab storage
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Shows current gitlab-storage-cleaner version

Flags:
  -h, --help               help for gitlab-storage-cleaner
  -l, --log-level string   set logging level

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
  -l, --log-level string   set logging level
```