#!/usr/bin/env bash
set -euo pipefail

# ─── Config ───────────────────────────────────────────────────────────────────
APP_NAME="stone"
OUTPUT_DIR="./dist"
MAIN_PKG="."                      # path to your main package
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}"
LDFLAGS="-s -w -X main.Version=${VERSION}"

# ─── Targets ──────────────────────────────────────────────────────────────────
# Format: "GOOS:GOARCH:output-suffix"
TARGETS=(
  "linux:amd64:"
  "linux:arm64:"
  "windows:amd64:.exe"
  "darwin:amd64:"
  "darwin:arm64:"
)

# ─── Helpers ──────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log()  { echo -e "${GREEN}[build]${NC} $*"; }
warn() { echo -e "${YELLOW}[warn]${NC}  $*"; }
die()  { echo -e "${RED}[error]${NC} $*" >&2; exit 1; }

require() {
  command -v "$1" &>/dev/null || die "'$1' is required but not installed."
}

# ─── Pre-flight ───────────────────────────────────────────────────────────────
require go
require zip

if [[ "${1:-}" == "clean" ]]; then
  log "Cleaning ${OUTPUT_DIR}/"
  rm -rf "${OUTPUT_DIR}"
  exit 0
fi

mkdir -p "${OUTPUT_DIR}"
log "Building ${APP_NAME} @ ${VERSION}"
echo ""

# ─── Build all targets ────────────────────────────────────────────────────────
for target in "${TARGETS[@]}"; do
  IFS=':' read -r os arch suffix <<< "$target"

  platform="${os}-${arch}"
  platform_dir="${OUTPUT_DIR}/${platform}"
  binary="${APP_NAME}${suffix}"
  zip_name="${APP_NAME}-${platform}-${VERSION}.zip"

  mkdir -p "${platform_dir}"

  printf "  %-14s %-8s -> %s/%s\n" "${os}" "${arch}" "${platform}" "${binary}"
  GOOS="${os}" GOARCH="${arch}" go build \
    -trimpath \
    -ldflags "${LDFLAGS}" \
    -o "${platform_dir}/${binary}" \
    "${MAIN_PKG}"

  # Zip the platform directory contents
  (cd "${platform_dir}" && zip -q "../${zip_name}" "${binary}")
  log "  zipped -> ${OUTPUT_DIR}/${zip_name}"
done

# ─── Universal macOS binary ───────────────────────────────────────────────────
echo ""
if command -v lipo &>/dev/null; then
  INTEL="${OUTPUT_DIR}/darwin-amd64/${APP_NAME}"
  ARM="${OUTPUT_DIR}/darwin-arm64/${APP_NAME}"
  UNIVERSAL_DIR="${OUTPUT_DIR}/darwin-universal"
  UNIVERSAL_ZIP="${OUTPUT_DIR}/${APP_NAME}-darwin-universal-${VERSION}.zip"

  if [[ -f "${INTEL}" && -f "${ARM}" ]]; then
    log "Creating universal macOS binary..."
    mkdir -p "${UNIVERSAL_DIR}"
    lipo -create -output "${UNIVERSAL_DIR}/${APP_NAME}" "${INTEL}" "${ARM}"
    (cd "${UNIVERSAL_DIR}" && zip -q "../${APP_NAME}-darwin-universal-${VERSION}.zip" "${APP_NAME}")
    log "  zipped -> ${UNIVERSAL_ZIP}"
  fi
else
  warn "lipo not found — skipping universal macOS binary (only available on macOS)"
fi

# ─── Checksums (zips only) ────────────────────────────────────────────────────
echo ""
log "Generating checksums..."
CHECKSUM_FILE="${OUTPUT_DIR}/checksums.sha256"

if command -v sha256sum &>/dev/null; then
  (cd "${OUTPUT_DIR}" && sha256sum ./*.zip > "checksums.sha256")
elif command -v shasum &>/dev/null; then
  (cd "${OUTPUT_DIR}" && shasum -a 256 ./*.zip > "checksums.sha256")
else
  warn "No sha256 tool found — skipping checksums"
fi

[[ -f "${CHECKSUM_FILE}" ]] && log "  -> ${CHECKSUM_FILE}"

# ─── Summary ──────────────────────────────────────────────────────────────────
echo ""
log "Done. Layout:"
find "${OUTPUT_DIR}" | sort | sed 's|[^/]*/|  |g'
