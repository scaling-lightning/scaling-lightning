#!/bin/bash
TAG=`git describe --tags --abbrev=0`
HASH=`git rev-parse --short HEAD`
DATE=`date`
echo "TAG=${TAG}"
echo "HASH=${HASH}"
echo "DATE=${DATE}"
echo -n "${TAG} | ${HASH} | ${DATE}" > version.txt
