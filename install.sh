#!/usr/bin/env sh
# install.sh — installs the latest course-builder binary
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/insigmo/course-builder/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/insigmo/course-builder/main/install.sh | sh -s -- --dir /usr/local/bin

set -e

REPO="insigmo/course-builder"
BINARY="course-builder"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Parse --dir flag
while [ $# -gt 0 ]; do
  case "$1" in
    --dir) INSTALL_DIR="$2"; shift 2 ;;
    *) shift ;;
  esac
done

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS. Use install.ps1 on Windows." >&2
    exit 1
    ;;
esac

# Get latest release tag
echo "Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine latest release." >&2
  exit 1
fi

FILENAME="${BINARY}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

TMP="$(mktemp)"
echo "Downloading ${FILENAME} (${LATEST})..."
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

mkdir -p "$INSTALL_DIR"
mv "$TMP" "${INSTALL_DIR}/${BINARY}"

echo ""
echo "✅ course-builder ${LATEST} installed to ${INSTALL_DIR}/${BINARY}"
echo "   Run: course-builder --help"
