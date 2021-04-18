#!/bin/bash
SCRIPT=$(readlink -f $0)
SCRIPTPATH=`dirname $SCRIPT`
PROJECT_DIR=`dirname $SCRIPTPATH`
PROJECT_DAT=`cat $PROJECT_DIR/project.json`
PROJECT_NAM=`echo $PROJECT_DAT | jq .Package -r`
PROJECT_VER=`echo $PROJECT_DAT | jq .Version -r`

echo
echo "Packaging $PROJECT_NAM v$PROJECT_VER"
echo

INSTALL_DIR="/usr/bin"
INSTALL_DIR=$(realpath -sm "$INSTALL_DIR")

echo "Using the following folders for install"
echo "PROJECT_DIR: $PROJECT_DIR"
echo "INSTALL_DIR: $INSTALL_DIR"

PACK="${PROJECT_DIR}/package/${PROJECT_NAM}_${PROJECT_VER}"
PACK_BIN=$(realpath -sm "$PACK/$INSTALL_DIR")
PACK_DEB=$(realpath -sm "$PACK/DEBIAN")
mkdir -p "$PACK"
mkdir -p "$PACK_BIN"
mkdir -p "$PACK_DEB"

# Copying files
echo "Copying binary files"
cp dist/mcversions "$PACK_BIN"

echo "Generating meta data"
PACK_CON=$(realpath -sm "$PACK_DEB/control")
echo "" > "$PACK_CON"
for row in $(echo "$PROJECT_DAT" | jq -r 'keys[]'); do
  value=`echo "$PROJECT_DAT" | jq -r ".[\"$row\"]"`
  echo "$row: $value" >> "$PACK_CON"
done

echo "Building package"
dpkg-deb --build "$PACK"
echo "Signing package"
dpkg-sig -k OnPointCoding --sign repo "$PACK.deb"
echo "Package complete:"
