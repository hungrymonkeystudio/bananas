#!/usr/bin/env bash 

set -euo pipefail

APP_NAME="bananas"
LOCAL_BIN="$HOME/.local/bin"
LOCAL_SHARE="$HOME/.local/share/$APP_NAME"

GLOBAL_BIN="/usr/local/bin"
GLOBAL_SHARE="/usr/local/share/$APP_NAME"

echo "Requires sudo for full uninstall"

echo "===Removing executable from $LOCAL_BIN==="
rm -f "$LOCAL_BIN/$APP_NAME"
echo "===Removing executable from $GLOBAL_BIN (requires sudo)==="
sudo rm -f "$GLOBAL_BIN/$APP_NAME"
echo "===Removing resources from $LOCAL_SHARE==="
rm -rf "$LOCAL_SHARE"
echo "===Removing resources from $GLOBAL_SHARE (requires sudo)==="
sudo rm -rf "$GLOBAL_SHARE"

echo "Uninstall successful"
