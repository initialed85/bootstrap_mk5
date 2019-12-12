#!/usr/bin/env bash

rm -fr dist >/dev/null 2>&1 || true
mkdir -p dist

go build -a -o dist/generate_mac cmd/generate_mac/main.go
go build -a -o dist/generate_ip cmd/generate_ip/main.go
go build -a -o dist/add_routes cmd/add_routes/main.go