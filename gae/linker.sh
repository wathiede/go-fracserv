#!/bin/sh
ln -vnsf $(python -c 'import os,sys;print os.path.realpath(sys.argv[1])' $1) github.com/wathiede/go-fracserv/$(basename $1)
