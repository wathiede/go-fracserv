#!/bin/bash

mkdir -p code.google.com/p/go-fracserv/
pushd code.google.com/p/
if [ -d "freetype-go" ]; then
    echo "Found freetype-go, updating"
    pushd freetype-go
    hg up
    popd
else
    echo "No checkout of freetype-go, fetching"
    hg clone http://code.google.com/p/freetype-go/
fi
popd

find ..  -maxdepth 1 -mindepth 1 -type d -not \( -name .git -o -name standalone -o -name gae -o -name lnp \) -exec ./linker.sh {} \;
ln -nsf ../fracserv/templates/ templates
# Pending a more long term solution from the SDK, remove problematic files:
rm code.google.com/p/freetype-go/example/round/main.go
