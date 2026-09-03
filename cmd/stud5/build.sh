# build.sh
APPNAME=$(basename "$PWD")

echo "Packaging $APPNAME..."
fyne package --release --icon=splash.png --id=com.scsi.$APPNAME

mkdir -p ../../dist
cp "$APPNAME.tar.xz" ../../dist/

echo "Done."
