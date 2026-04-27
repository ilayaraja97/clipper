#!/bin/bash

BINNAME="${BINNAME:-clipper}"
BINDIR="${BINDIR:-${HOME}/.local/bin}"

echo "Uninstallation of Clipper ..."
echo

rm $BINDIR/$BINNAME
rm "${HOME}/.config/${BINNAME}.json"

echo
echo "Uninstallation of Clipper complete!"
