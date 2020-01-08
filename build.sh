#!/usr/bin/env bash

set -x
set -e

echo "cleaning..."
rm -fr dist >/dev/null 2>&1 || true
mkdir -p dist
echo ""

echo "building tools to be deployed..."
go build -v -o dist/generate_mac cmd/generate_mac/main.go
go build -v -o dist/generate_ip cmd/generate_ip/main.go
go build -v -o dist/add_routes cmd/add_routes/main.go
go build -v -o dist/castinator cmd/castinator/main.go
echo ""

echo "hashing and stamping..."
tar -c dist | sha1sum >dist/build_hash.txt
date >dist/build_date.txt
echo ""

echo "building tools to assist deployment..."
unset GOOS
unset GOARCH
rm -fr deploy/find_mk5s 2>&1 || true
go build -v -o deploy/find_mk5s cmd/find_mk5s/main.go
echo ""
