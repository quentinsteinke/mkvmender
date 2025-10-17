#!/bin/sh
# MKV Mender Installation Script
# Supports: Linux (x86_64, ARM64), macOS (Intel, Apple Silicon)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="quentinsteinke/mkvmender"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="mkvmender"

# Helper functions
info() {
    printf "${BLUE}==>${NC} %s\n" "$1" >&2
}

success() {
    printf "${GREEN}✓${NC} %s\n" "$1" >&2
}

error() {
    printf "${RED}✗${NC} %s\n" "$1" >&2
    exit 1
}

warn() {
    printf "${YELLOW}!${NC} %s\n" "$1" >&2
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        *)          error "Unsupported operating system: $(uname -s)";;
    esac
}

# Detect architecture
detect_arch() {
    local arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64)   echo "amd64";;
        aarch64|arm64)  echo "arm64";;
        *)              error "Unsupported architecture: $arch";;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    info "Fetching latest version..."

    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        error "Either curl or wget is required to download mkvmender"
    fi

    if [ -z "$LATEST_VERSION" ]; then
        error "Failed to get latest version"
    fi

    success "Latest version: $LATEST_VERSION"
}

# Download binary
download_binary() {
    local os="$1"
    local arch="$2"
    local version="$3"

    local filename="mkvmender-${os}-${arch}"
    if [ "$os" = "darwin" ] && [ "$arch" = "amd64" ]; then
        # macOS Intel build might have different naming
        filename="mkvmender-${os}-${arch}"
    fi

    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${filename}"
    local tmp_file="/tmp/${filename}"

    info "Downloading mkvmender from $download_url..."

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$tmp_file" || error "Failed to download binary"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$tmp_file" || error "Failed to download binary"
    fi

    success "Downloaded successfully"
    echo "$tmp_file"
}

# Install binary
install_binary() {
    local tmp_file="$1"
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"

    info "Installing to $install_path..."

    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"

    # Move binary to install location
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_file" "$install_path"
    else
        # Need sudo for system-wide installation
        warn "Requesting elevated privileges to install to $INSTALL_DIR"
        sudo mv "$tmp_file" "$install_path"
    fi

    # Make executable
    chmod +x "$install_path" 2>/dev/null || sudo chmod +x "$install_path"

    success "Installed to $install_path"
}

# Check if directory is in PATH
check_path() {
    local dir="$1"

    if ! echo "$PATH" | grep -q "$dir"; then
        warn "$dir is not in your PATH"
        info "Add it to your PATH by adding this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        printf "\n    export PATH=\"\$PATH:%s\"\n\n" "$dir" >&2
        info "Then reload your shell with: source ~/.bashrc (or ~/.zshrc)"
    fi
}

# Main installation
main() {
    echo ""
    info "MKV Mender Installer"
    echo ""

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: $OS $ARCH"

    # Get latest version
    get_latest_version

    # Download binary
    TMP_FILE=$(download_binary "$OS" "$ARCH" "$LATEST_VERSION")

    # Install binary
    install_binary "$TMP_FILE"

    # Check PATH
    check_path "$INSTALL_DIR"

    echo ""
    success "MKV Mender installed successfully!"
    echo ""
    info "Get started with:"
    printf "    %s register\n" "$BINARY_NAME" >&2
    echo ""
    info "For more information:"
    printf "    %s --help\n" "$BINARY_NAME" >&2
    printf "    https://github.com/%s\n" "$GITHUB_REPO" >&2
    echo ""
}

# Run main installation
main
