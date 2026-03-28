#!/bin/bash

BINNAME="${BINNAME:-clipper}"
BINDIR="${BINDIR:-/usr/local/bin}"

echo "Uninstallation of Clipper ..."
echo

sudo rm $BINDIR/$BINNAME
sudo rm "${HOME}/.config/${BINNAME}.json"

echo
echo "Uninstallation of Clipper complete!"
