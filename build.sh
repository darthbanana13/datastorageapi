#!/usr/bin/env bash

if [[ ! -x "$(which go)" ]]; then
    echo "Please install Go version 1.17 or higher"
    exit 1
fi

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
go build -o "${SCRIPTPATH}/dataApi" cmd/server/main.go
