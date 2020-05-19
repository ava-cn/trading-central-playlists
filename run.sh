#! /bin/sh
set -e

# Go Binary Path
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

# make build_linux_amd64

# rebuild image
# docker build -t curder/trading-central-playlists .

docker-compose build

# sudo systemctl start docker-compose-trading-central-playlists
