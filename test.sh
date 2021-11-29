#!/usr/bin/env bash

if [[ ! -x "$(which go)" ]]; then
    echo "Error! Go executable not found in system path. Please install Go version 1.17 or higher"
    exit 1
fi

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

cd "${SCRIPTPATH}"

go test internal/aqlBuilder/aqlBuilder_test.go
go test internal/service/chatFilterBuilderDecorator/chatFilterBuilderDecorator_test.go
