#!/bin/bash
set -x
set -e
cd ../cmd/crosspostcontrol
go build
mv crosspostcontrol ../../build
cd ../../build
tar cvfz crosspostcontrol-plugin.tar.gz plugin.yaml crosspostcontrol
