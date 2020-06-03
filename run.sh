#!/usr/bin/env bash

if ["$(uname)"=="Darwin"];then

# Mac OS X 操作系统
export GOROOT=/usr/local/go # Go Binary Path

elif ["$(expr substr $(uname -s) 1 5)"=="Linux"];then

# GNU/Linux操作系统
export GOROOT=/usr/local/Cellar/go/1.14/libexec # Go Binary Path

fi


export PATH=$PATH:$GOROOT/bin

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

make build_linux_amd64

# rebuild image
# docker build -t curder/trading-central-playlists .

docker-compose build

docker-compose up -d
