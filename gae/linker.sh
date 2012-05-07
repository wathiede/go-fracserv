#!/bin/sh
ln -nsf $(readlink -f $1) code.google.com/p/go-fracserv/$(basename $1)
