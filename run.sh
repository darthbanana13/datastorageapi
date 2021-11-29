#!/usr/bin/env bash

if [[ ! -x "$(which go)" ]]; then
    echo "Error! Go executable not found in system path. Please install Go version 1.17 or higher"
    exit 1
fi

if [[ ! -x "$(which docker-compose)" ]]; then
    echo "Error! Missing executable docker-compose. Please install it"
    exit 1
fi

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

if [[ ! -f "${SCRIPTPATH}/.env" ]]; then
    echo ".env does not exist, please copy .env.example into .env and change the values accordingly"
    exit 1
fi

EXE=(go run cmd/server/main.go)
if [[ -f "${SCRIPTPATH}/dataApi" ]]; then
    EXE=("./dataApi")
fi

cd $SCRIPTPATH
docker-compose up -d
"${EXE[@]}"
