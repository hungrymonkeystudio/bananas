#!/usr/bin/env bash 
    
set -euo pipefail

APP_NAME="bananas"
LOCAL_BIN="$HOME/.local/bin"
LOCAL_SHARE="$HOME/.local/share/$APP_NAME"
# LOCAL_CONFIG="$HOME/.config/$APP_NAME"
# LOCAL_CACHE="$HOME/.cache/$APP_NAME"

GLOBAL_BIN="/usr/local/bin"
GLOBAL_SHARE="/usr/local/share/$APP_NAME"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="${SCRIPT_DIR%/*}"
RESOURCE_DIR="$PARENT_DIR/resources"

INSTALL_MODE="local"

if [[ "${1:-}" == "--global" ]]; then
    INSTALL_MODE="global"
fi

echo "===Checking Go is installed==="
if ! command -v go >/dev/null 2>&1; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

echo "===Checking Go version==="
VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?')
if [[ "$VERSION" != "1.24.6" ]]; then
    echo "WARNING: Program might not be compatible on $VERSION"
    echo "Can largely ignore since will probably be compatible"
fi

echo "===Building executable==="
cd $PARENT_DIR
go build . 

if [[ "$INSTALL_MODE" == "local" ]]; then
    echo "===Installing locally==="
    mkdir -p "$LOCAL_BIN" "$LOCAL_SHARE"
    mv "$APP_NAME" "$LOCAL_BIN/"
    if [[ -d "$RESOURCE_DIR" ]]; then
        cp -r "$RESOURCE_DIR"/* "$LOCAL_SHARE/"
    fi
    if [[ ":$PATH:" != *":$LOCAL_BIN:"* ]]; then
        echo "WARNING: $LOCAL_BIN is not in your PATH"
        echo "Add this line to your shell config"
        echo "  export PATH=\$PATH:$LOCAL_BIN"
    fi
    echo "Installed $APP_NAME to $LOCAL_BIN"
    echo "Resources in $LOCAL_SHARE"
else
    echo "===Installing globally (requires sudo)==="
    echo "===Installing globally (requires sudo)==="
    sudo mkdir -p "$GLOBAL_BIN" "$GLOBAL_SHARE"
    sudo mv "$APP_NAME" "$GLOBAL_BIN/"
    if [[ -d "$RESOURCE_DIR" ]]; then
        cp -r "$RESOURCE_DIR"/* "$GLOBAL_SHARE/"
    fi
    echo "Installed $APP_NAME to $GLOBAL_BIN"
    echo "Resources in $GLOBAL_SHARE"
fi
