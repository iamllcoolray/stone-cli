#!/usr/bin/env bash
set -euo pipefail

REPO="iamllcoolray/stone-cli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# ─── Detect OS + Arch ────────────────────────────────────────────────────────
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)  ;;
  darwin) ;;
  *) echo "Unsupported OS: $OS. For Windows, see https://github.com/${REPO}#installation"; exit 1 ;;
esac

case "$ARCH" in
  x86_64)        ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# macOS: report which Mac variant was detected
if [[ "$OS" == "darwin" ]]; then
  if [[ "$ARCH" == "amd64" ]]; then
    echo "Detected: macOS Intel (x86_64)"
  else
    echo "Detected: macOS Apple Silicon (arm64)"
  fi
fi

# ─── Download ─────────────────────────────────────────────────────────────────
# get latest version tag from GitHub API
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

ZIP="stone-${OS}-${ARCH}-${VERSION}.zip"
URL="https://github.com/${REPO}/releases/latest/download/${ZIP}"
TMP_DIR=$(mktemp -d)

echo "Downloading ${ZIP}..."
curl -fsSL "$URL" -o "${TMP_DIR}/${ZIP}"

# ─── Extract ──────────────────────────────────────────────────────────────────
echo "Extracting..."
unzip -q "${TMP_DIR}/${ZIP}" -d "${TMP_DIR}/stone-install"

# ─── Install ──────────────────────────────────────────────────────────────────
# Try system-wide first, fall back to ~/.local/bin if no permission
if [ -w "$INSTALL_DIR" ] || sudo -n true 2>/dev/null; then
  echo "Installing to ${INSTALL_DIR}..."
  sudo mv "${TMP_DIR}/stone-install/stone" "${INSTALL_DIR}/stone"
  sudo chmod +x "${INSTALL_DIR}/stone"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "No sudo access — installing to ${INSTALL_DIR}..."
  mv "${TMP_DIR}/stone-install/stone" "${INSTALL_DIR}/stone"
  chmod +x "${INSTALL_DIR}/stone"

  # Check if ~/.local/bin is on PATH
  if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "Adding ${INSTALL_DIR} to PATH..."

    SHELL_NAME=$(basename "${SHELL:-/bin/bash}")
    case "$SHELL_NAME" in
      zsh)  RC="$HOME/.zshrc" ;;
      fish) RC="$HOME/.config/fish/config.fish" ;;
      *)    RC="$HOME/.bashrc" ;;
    esac

    echo "" >> "$RC"
    echo "# Added by stone installer" >> "$RC"
    echo "export PATH=\"\$PATH:${INSTALL_DIR}\"" >> "$RC"
    echo "Added to ${RC} — run: source ${RC}"
  fi
fi

# ─── Cleanup ──────────────────────────────────────────────────────────────────
rm -rf "$TMP_DIR"

# ─── Done ─────────────────────────────────────────────────────────────────────
echo ""
echo "stone installed successfully."
echo "Run: stone --version"