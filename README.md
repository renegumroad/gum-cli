# gum-cli

## What is gum?

`gum` is our in-house developer toolkit!

## Installation

Download the latest release from GitHub, using your own PAT (Personal Access Token). Replace `<OS>` by `linux`/`darwin` and `<ARCH>` by `arm64`/`amd64` where appropriate.

```shell
URL=$(curl -s -H "Authorization: token <TOKEN>" https://api.github.com/repos/renegumroad/gum-cli/releases/latest | jq -r '.assets[] | select(.name == "gum_<OS>_<ARCH>") | .url')
curl -Lv -J -H "Accept: application/octet-stream" -H "Authorization: token <TOKEN>" -o gum $URL
chmod +x gum
./gum --help
```

You should see the top-level help information after the above commands

## Usage

## `gum init`

Configures workstation with prerequisites for gum

## `gum dev up`

Configures development dependencies in `gum.yml`

```yaml
# gum.y(a)ml

up:
  - action: golang
  - brew:
      - name: jq
      - name: yq
```

### Logging

Logging can be tweaked via `--log-level=<level>` flag.

**Example using debug level:**

```shell
gum <subcommand> --log-level=debug
```
