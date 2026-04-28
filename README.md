# stone 🪨

A cross-platform CLI tool for checking and updating [utiliti](https://gurkenlabs.itch.io/litiengine) distributed on [itch.io](https://itch.io).

Named after *lithos* — the Greek root behind [LITIengine](https://litiengine.com).

---

## Features

- Check for new versions of any public itch.io project
- Download and extract the latest release for your platform
- Replace the existing installation in place
- Cross-platform: Linux, macOS (Intel + Apple Silicon), Windows

---

## Installation

**Linux / macOS:**
```bash
curl -L https://github.com/iamllcoolray/stone-cli/releases/latest/download/install.sh | bash
```

**Windows:**

Download the latest zip from the [Releases](https://github.com/iamllcoolray/stone-cli/releases/latest) page, extract `stone.exe`, and add its location to your system `PATH`.

| Platform         | File                              |
|------------------|-----------------------------------|
| Linux (x86_64)   | `stone-linux-amd64.zip`           |
| Linux (ARM64)    | `stone-linux-arm64.zip`           |
| macOS (Intel)    | `stone-darwin-amd64.zip`          |
| macOS (M1/M2/M3) | `stone-darwin-arm64.zip`          |
| Windows (x86_64) | `stone-windows-amd64.zip`         |

### Custom install directory

```bash
INSTALL_DIR=~/.local/bin curl -L https://github.com/iamllcoolray/stone-cli/releases/latest/download/install.sh | bash
```

---

## Configuration

On first run, stone looks for a config file at:

| Platform | Path |
|----------|------|
| Linux / macOS | `~/.config/stone/config.toml` |
| Windows | `%APPDATA%\stone\config.toml` |

Create the file manually or run `stone init` to generate it interactively.

**`config.toml`:**
```toml
install_path = "/path/to/your/game"
last_version = "0.11.1"
```

---

## Usage

```bash
# Set up stone and install utiLITI
stone init

# Veiw the local install path
stone config

# Configure the itch.io API key
stone config --api-key

# Configure the local install path
stone config --install-path

# Check if a new version is available
stone check

# Download and install the latest version
stone update

# Force update even if already on the latest version
stone update --force

# Remove utiLITI from the system
stone remove

# Remove utiLITI and stone from the system
stone remove --all

# Print the current stone version
stone version
```

---

## Verifying Downloads

Each release includes a `checksums.sha256` file. To verify your download:

**Linux / macOS:**
```bash
sha256sum -c checksums.sha256 --ignore-missing
```

**Windows (PowerShell):**
```powershell
Get-FileHash stone-windows-amd64.zip -Algorithm SHA256
```

Compare the output against the matching line in `checksums.sha256`.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

MIT