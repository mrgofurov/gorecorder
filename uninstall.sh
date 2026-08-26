#!/usr/bin/env bash
set -e

echo "🗑️ Uninstalling GoRecorder..."

rm -rf "$HOME/.local/share/gorecorder"
rm -f "$HOME/.local/bin/gorecorder"
rm -f "$HOME/.local/share/applications/gorecorder.desktop"
rm -f "$HOME/.local/share/icons/hicolor/512x512/apps/gorecorder.png"

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$HOME/.local/share/applications" 2>/dev/null || true
fi

echo "✅ GoRecorder has been uninstalled."
