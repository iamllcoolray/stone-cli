# Contributing to stone

Thanks for your interest in contributing. This document covers everything you need to get up and running.

For a deeper look at the codebase internals, architecture decisions, and release process see [DEVELOPER.md](DEVELOPER.md).

---

## Getting Started

**Requirements:**
- Go 1.21+
- `zip` (for the build script)
- Git

```bash
git clone https://github.com/iamllcoolray/stone-cli
cd stone-cli
go mod download
```

Build for your current platform:

```bash
go build -o stone .
./stone --version
```

---

## Project Structure

```
stone-cli/
├── .github/workflows/
│   └── release.yml       # CI release workflow
├── cmd/
│   └── root.go           # cobra root + subcommand registration
├── internal/
│   ├── api/
│   │   └── api.go        # itch.io api, version fetch
│   ├── updater/
│   │   └── updater.go    # download, extract, replace logic
│   └── config/
│       └── config.go     # load/save config (viper)
├── main.go               # entrypoint
├── go.mod
├── go.sum
├── build.sh              # cross-platform build script
├── install.sh            # curl installer
└── .gitignore
```

---

## Workflow

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Run `go vet ./...` and `go test ./...`
4. Open a pull request with a clear description of what changed and why

Branch naming:
- `feat/short-description` for new features
- `fix/short-description` for bug fixes
- `chore/short-description` for maintenance

---

## Commit Style

Use short, lowercase imperative commit messages:

```
add darwin amd64 and arm64 build targets
fix path detection on windows
update install dir fallback logic
```

No ticket numbers, no emoji, no periods at the end.

---

## Reporting Issues

Open an issue on [GitHub](https://github.com/iamllcoolray/stone-cli/issues) with:
- Your OS and architecture
- The command you ran
- The full output or error message
- stone version (`stone --version`)

---

## License

By contributing you agree that your contributions will be licensed under the MIT license.