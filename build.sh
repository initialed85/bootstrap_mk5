#!/usr/bin/env bash

set -x
set -e

rm -fr dist >/dev/null 2>&1 || true
mkdir -p dist

go build -x -v -o dist/generate_mac cmd/generate_mac/main.go
go build -x -v -o dist/generate_ip cmd/generate_ip/main.go
go build -x -v -o dist/add_routes cmd/add_routes/main.go
go build -x -v -o dist/castinator cmd/castinator/main.go
