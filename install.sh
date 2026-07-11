#!/usr/bin/env bash
set -euo pipefail

GITHUB_REPO="neracu/vibe-shield"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="vibe-shield"

log()  { printf '\033[1;36m==>\033[0m %s\n' "$*"; }
ok()   { printf '\033[1;32m✔\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m!\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31m✖\033[0m %s\n' "$*" >&2; exit 1; }

detect_platform() {
  local raw_os raw_arch

  log "Checking OS and architecture..."
  raw_os="$(uname -s)"
  raw_arch="$(uname -m)"

  case "$raw_os" in
    Darwin)
      os="darwin"
      case "$raw_arch" in
        x86_64)  arch="amd64"; label="macOS (Intel)" ;;
        arm64)   arch="arm64"; label="macOS (Apple Silicon / M1/M2/M3)" ;;
        *) fail "Unsupported macOS architecture: $raw_arch" ;;
      esac
      ;;
    Linux)
      os="linux"
      case "$raw_arch" in
        x86_64|amd64) arch="amd64"; label="Linux (amd64)" ;;
        aarch64|arm64) arch="arm64"; label="Linux (arm64)" ;;
        *) fail "Unsupported Linux architecture: $raw_arch" ;;
      esac
      ;;
    *)
      fail "Unsupported operating system: $raw_os (this installer supports macOS and Linux only)"
      ;;
  esac

  asset_name="vibe-shield-${os}-${arch}"
  download_url="https://github.com/${GITHUB_REPO}/releases/latest/download/${asset_name}"
  ok "Detected ${label}"
}

download_binary() {
  local tmp_dir tmp_file

  tmp_dir="$(mktemp -d)"
  tmp_file="${tmp_dir}/${BINARY_NAME}"

  trap 'rm -rf "$tmp_dir"' EXIT

  log "Downloading Vibe-Shield from GitHub Releases..."
  log "URL: ${download_url}"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$download_url" -o "$tmp_file"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$download_url" -O "$tmp_file"
  else
    fail "Neither curl nor wget found. Install one of them and retry."
  fi

  chmod +x "$tmp_file"
  ok "Download complete"
  INSTALL_SRC="$tmp_file"
}

install_binary() {
  local target="${INSTALL_DIR}/${BINARY_NAME}"

  log "Installing to ${target}..."

  if [[ -w "$INSTALL_DIR" ]]; then
    mv "$INSTALL_SRC" "$target"
  elif command -v sudo >/dev/null 2>&1; then
    warn "${INSTALL_DIR} is not writable — sudo required"
    sudo mv "$INSTALL_SRC" "$target"
  else
    fail "Cannot write to ${INSTALL_DIR}. Re-run with sudo or fix permissions."
  fi

  ok "Successfully installed!"
  log "Run: vibe-shield <your-command>"
}

main() {
  printf '\n\033[1;35mVibe-Shield Installer\033[0m\n\n'
  detect_platform
  download_binary
  install_binary
  printf '\n'
}

main "$@"
