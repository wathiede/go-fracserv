#!/bin/sh
ln -vnsf $(python -c 'import os,sys;print os.path.realpath(sys.argv[1])' $1) code.google.com/p/go-fracserv/$(basename $1)
