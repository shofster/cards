#!/bin/sh

APP="$1"
ARCH="$2"

# ------------------------------------------------------------
# Usage / defaults
# ------------------------------------------------------------
[ -z "$APP" ] && { echo "Usage: ./games.sh <app|all>"; exit 1; }

DIST="dist"
OS=$(uname)

# ------------------------------------------------------------
# Build ALL apps
# ------------------------------------------------------------
if [ "$APP" = "all" ]; then
    for DIR in cmd/*; do
        [ -d "$DIR" ] || continue
        NAME=$(basename "$DIR")
        ./games.sh "$NAME" "$ARCH"
    done
    exit 0
fi

# ------------------------------------------------------------
# Build ONE app
# ------------------------------------------------------------
echo
echo "=== Packaging $APP for $OS/$ARCH ==="

cd "cmd/$APP" || { echo "ERROR: cmd/$APP not found"; exit 1; }
fyne package --release --icon=splash.png --id="com.scsi.$APP"
cd - >/dev/null

mkdir -p "$DIST"

# ------------------------------------------------------------
# macOS: .app + .dmg
# ------------------------------------------------------------
if [ "$OS" = "Darwin" ]; then
    if [ -d "cmd/$APP/$APP.app" ]; then
        chmod +x "cmd/$APP/$APP.app/Contents/MacOS/$APP"
        cp -R "cmd/$APP/$APP.app" "$DIST/"

        echo "Creating DMG for $APP"
        hdiutil create -volname "$APP" \
            -srcfolder "cmd/$APP/$APP.app" \
            -ov -format UDZO "$DIST/$APP.dmg"
    else
        echo "ERROR: No .app found for macOS"
        exit 1
    fi

# ------------------------------------------------------------
# Linux: .tar.xz
# ------------------------------------------------------------
else
    if [ -f "cmd/$APP/$APP.tar.xz" ]; then
        cp "cmd/$APP/$APP.tar.xz" "$DIST/"
    else
        echo "ERROR: No .tar.xz found for Linux"
        exit 1
    fi
fi

echo "Done: $APP → $DIST"
