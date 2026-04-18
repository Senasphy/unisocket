#!/usr/bin/env sh
set -eu

REPO="${UNISOCKET_REPO:-Senasphy/unisocket}"
BIN_NAME="${UNISOCKET_BIN:-unisocket}"
VERSION="${1:-${UNISOCKET_VERSION:-}}"
INSTALL_DIR="${UNISOCKET_INSTALL_DIR:-/usr/local/bin}"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

need_cmd curl
need_cmd tar
need_cmd uname

detect_os() {
  case "$(uname -s)" in
    Linux) echo "linux" ;;
    Darwin) echo "darwin" ;;
    *) echo "unsupported OS"; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported arch"; exit 1 ;;
  esac
}

resolve_latest_version() {
  curl -fsSL -H "User-Agent: unisocket-installer" \
    "https://api.github.com/repos/${REPO}/releases/latest" \
  | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' \
  | head -n 1
}

OS="$(detect_os)"
ARCH="$(detect_arch)"

[ -z "${VERSION}" ] && VERSION="$(resolve_latest_version)"
[ -z "${VERSION}" ] && { echo "failed to resolve version"; exit 1; }

VERSION="${VERSION#v}"

ARTIFACT="${BIN_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARTIFACT}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

ARCHIVE="$TMP_DIR/$ARTIFACT"

echo "Downloading $URL ..."
curl -fL "$URL" -o "$ARCHIVE"

tar -xzf "$ARCHIVE" -C "$TMP_DIR"

SRC="$TMP_DIR/$BIN_NAME"
DEST="${INSTALL_DIR%/}/$BIN_NAME"

[ -f "$SRC" ] || { echo "binary not found"; exit 1; }

mkdir -p "$INSTALL_DIR" 2>/dev/null || true

if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "$SRC" "$DEST"
else
  sudo install -m 0755 "$SRC" "$DEST"
fi

echo "installed $BIN_NAME $VERSION → $DEST"
