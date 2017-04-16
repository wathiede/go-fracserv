#!/bin/bash

mkdir -p github.com/wathiede/go-fracserv/
pushd github.com/wathiede/
if [ -d "freetype-go" ]; then
    echo "Found freetype-go, updating"
    pushd freetype-go
    hg up
    popd
else
    echo "No checkout of freetype-go, fetching"
    hg clone http://github.com/wathiede/freetype-go/
fi
popd

find ..  -maxdepth 1 -mindepth 1 -type d -not \( -name .git -o -name standalone -o -name gae -o -name lnp \) -exec ./linker.sh {} \;
ln -nsf ../fracserv/templates/ templates
# Pending a more long term solution from the SDK, remove problematic files:
rm -f github.com/wathiede/freetype-go/example/round/main.go
