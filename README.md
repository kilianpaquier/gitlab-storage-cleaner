# gitlab-storage-cleaner <!-- omit in toc -->

- [How to use ?](#how-to-use-)
  - [Easy install](#easy-install)
  - [Specific install](#specific-install)
- [Commands](#commands)

## How to use ?

### Easy install

```sh
go install github.com/kilianpaquier/gitlab-storage-cleaner/cmd/gitlab-storage-cleaner@latest
```

### Specific install

You can either download a specific release asset at https://github.com/kilianpaquier/gitlab-storage-cleaner/-/releases or run the below shell commands to directly download an asset (just change the `BINARY` value).

```sh
BINARY="gitlab-storage-cleaner_linux_amd64"
DEST="gitlab-storage-cleaner"
DOWNLOAD_URL="$(curl -fsSL https://gitlab.com/api/v4/projects/54329665/releases/permalink/latest | jq -r ".assets.links[] | select(.name == \"$BINARY\") | .direct_asset_url")"
curl -fsSL "$DOWNLOAD_URL" -o "$DEST"
chmod +x "$DEST"
```

## Commands