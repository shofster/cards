# build.sh  (NO shebang)
APPNAME=$(basename "$PWD")

echo "Packaging $APPNAME..."
fyne package --release --icon=splash.png --id=com.scsi.$APPNAME

mkdir -p ../../dist/arm
cp "$APPNAME.tar.xz" ../../dist/arm/

echo "Done."
