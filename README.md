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
 
## itch.io API Key
 
stone uses the [itch.io API](https://itch.io/docs/api/overview) to check for new versions of utiLITI and retrieve download URLs. An API key is required.
 
### Getting a key
 
1. Log in to [itch.io](https://itch.io)
2. Go to [Settings → API keys](https://itch.io/user/settings/api-keys)
3. Click **Generate new API key**
4. Copy the key — you'll need it during `stone init`
### Important
 
stone requires a **personal API key** generated from your itch.io account settings. OAuth keys will not work for downloading files — this is an itch.io restriction, not a stone limitation.
 
Your API key is stored locally in your config file and is never transmitted anywhere other than directly to the itch.io API.
 
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
api_key      = "your_itchio_api_key"
install_path = "/path/to/utiliti"
last_version = "0.11.1"
```

`last_version` is written automatically by stone after a successful install or update — you do not need to set it manually.

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

# Create a new LITIengine project
stone new

# Check if a new version is available
stone check

# Download and install the latest version
stone update

# Force update even if already on the latest version
stone update --force

# Remove utiLITI from the system
stone remove

# Remove stone from the system
stone remove --stone

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