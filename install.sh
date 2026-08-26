#!/usr/bin/env bash
set -e

echo "🚀 Installing GoRecorder Desktop Application..."

INSTALL_DIR="$HOME/.local/share/gorecorder"
BIN_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
PIXMAP_DIR="$HOME/.local/share/pixmaps"

mkdir -p "$INSTALL_DIR"
mkdir -p "$BIN_DIR"
mkdir -p "$DESKTOP_DIR"
mkdir -p "$PIXMAP_DIR"
mkdir -p "$HOME/.local/share/icons/hicolor"/{48x48,64x64,128x128,256x256,512x512}/apps

# 1. Copy binary
cp -f "$(pwd)/build/bin/gorecorder" "$INSTALL_DIR/gorecorder"
chmod +x "$INSTALL_DIR/gorecorder"

# 2. Symlink to user PATH
ln -sf "$INSTALL_DIR/gorecorder" "$BIN_DIR/gorecorder"

# 3. Copy Icons in all standard resolutions
if [ -f "$(pwd)/build/appicon.png" ]; then
    cp -f "$(pwd)/build/appicon.png" "$INSTALL_DIR/icon.png"
    cp -f "$(pwd)/build/appicon.png" "$PIXMAP_DIR/gorecorder.png"
    cp -f "$(pwd)/build/appicon.png" "$HOME/.local/share/icons/hicolor/512x512/apps/gorecorder.png"
    ffmpeg -y -i "$(pwd)/build/appicon.png" -vf "scale=256:256" "$HOME/.local/share/icons/hicolor/256x256/apps/gorecorder.png" 2>/dev/null || true
    ffmpeg -y -i "$(pwd)/build/appicon.png" -vf "scale=128:128" "$HOME/.local/share/icons/hicolor/128x128/apps/gorecorder.png" 2>/dev/null || true
    ffmpeg -y -i "$(pwd)/build/appicon.png" -vf "scale=64:64" "$HOME/.local/share/icons/hicolor/64x64/apps/gorecorder.png" 2>/dev/null || true
    ffmpeg -y -i "$(pwd)/build/appicon.png" -vf "scale=48:48" "$HOME/.local/share/icons/hicolor/48x48/apps/gorecorder.png" 2>/dev/null || true
fi

# 4. Create .desktop launcher
cat <<EOF > "$DESKTOP_DIR/gorecorder.desktop"
[Desktop Entry]
Name=GoRecorder
Comment=High-Performance Screen Recorder for Ubuntu Wayland
Exec=$INSTALL_DIR/gorecorder
Icon=gorecorder
Terminal=false
Type=Application
Categories=AudioVideo;Recorder;Utility;
Keywords=screen;recorder;video;capture;gorecorder;wayland;
StartupWMClass=gorecorder
EOF

chmod +x "$DESKTOP_DIR/gorecorder.desktop"

# 5. Update desktop database and icon caches
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -f -t "$HOME/.local/share/icons/hicolor" 2>/dev/null || true
fi

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
fi

touch "$DESKTOP_DIR/gorecorder.desktop"

echo "✅ GoRecorder successfully installed!"
