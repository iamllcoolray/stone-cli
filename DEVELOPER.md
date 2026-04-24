# Developer Guide

Internal documentation covering architecture, build system, and release process for stone.

---

## Architecture

stone is structured as a standard Go CLI using cobra for commands and viper for config. All implementation logic lives under `internal/` so it cannot be imported outside the module.

```
cmd/          → CLI layer (cobra commands, flag definitions)
internal/     → all implementation, not importable externally
  api/        → itch.io HTTP scraping, version resolution
  updater/    → download, extract, file replacement
  config/     → config file read/write via viper
main.go       → entrypoint, Version var injected at build time
```

### Command flow

```
stone check
  └── config.Load()
  └── api.FetchLatest(gameID)
  └── compare against stored version
  └── print result

stone update
  └── config.Load()
  └── api.FetchLatest(gameID)
  └── updater.Download(url)
  └── updater.Extract(zipPath, installPath)
  └── config.SaveVersion(latest)
```

---

## Version Injection

The `Version` variable in `main.go` is set to `"dev"` by default and overwritten at build time via `-ldflags`:

```go
// main.go
var Version = "dev"
```

```bash
go build -ldflags "-X main.Version=1.2.3" -o stone .
```

The build script handles this automatically using `git describe`:

```bash
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
```

---

## Build System

All builds go through `build.sh`. It produces per-platform directories and zips under `dist/`:

```
dist/
├── linux-amd64/
│   └── stone
├── linux-arm64/
│   └── stone
├── darwin-amd64/
│   └── stone
├── darwin-arm64/
│   └── stone
├── windows-amd64/
│   └── stone.exe
├── stone-linux-amd64-v1.0.0.zip
├── stone-linux-arm64-v1.0.0.zip
├── stone-darwin-amd64-v1.0.0.zip
├── stone-darwin-arm64-v1.0.0.zip
├── stone-windows-amd64-v1.0.0.zip
└── checksums.sha256
```

Common commands:

```bash
./build.sh              # build all platforms
./build.sh clean        # wipe dist/
VERSION=1.2.3 ./build.sh  # explicit version override
```

`dist/` is gitignored — never commit it.

---

## Release Process

Releases are fully automated via `.github/workflows/release.yml`. The workflow triggers on any tag matching `v*.*.*`.

To cut a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions will:
1. Check out the repo with full tag history
2. Set up Go from the version in `go.mod`
3. Run `./build.sh`
4. Create a GitHub Release named after the tag
5. Upload all `dist/*.zip`, `dist/checksums.sha256`, and `install.sh` as release assets
6. Auto-generate release notes from commits since the last tag

The release assets are then immediately available at:
```
https://github.com/iamllcoolray/stone-cli/releases/latest/download/<filename>
```

---

## Config File

stone uses TOML via viper. The config file is resolved at runtime from platform-specific locations:

| Platform | Path |
|----------|------|
| Linux / macOS | `~/.config/stone/config.toml` |
| Windows | `%APPDATA%\stone\config.toml` |

**Fields:**

| Key | Type | Description |
|-----|------|-------------|
| `game_id` | string | The itch.io game ID or slug to track |
| `install_path` | string | Absolute path to the local installation |
| `last_version` | string | Stored after a successful update, used for comparison |

---

## itch.io Scraping

stone scrapes public itch.io pages rather than using the API — no key required. The scraper in `internal/api/itchio.go` hits the game's public page and parses the upload metadata to find the latest build for the current platform.

Platform detection maps `runtime.GOOS` to itch.io's platform tags:

| `runtime.GOOS` | itch.io tag |
|----------------|-------------|
| `linux` | `p_linux` |
| `windows` | `p_windows` |
| `darwin` | `p_osx` |

**Note:** itch.io's page structure can change without notice. If the scraper breaks after an itch.io update, the HTML selectors in `itchio.go` will need updating.

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI commands and flags |
| `github.com/spf13/viper` | Config file management |
| `archive/zip` | stdlib — zip extraction |
| `net/http` | stdlib — HTTP requests and downloads |