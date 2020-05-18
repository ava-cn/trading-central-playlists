#! /bin/sh
set -e

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

go build -v -o dist/trading-central-playlists

docker build -t curder/trading-central-playlists .

docker-compose up -d