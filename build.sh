#!/bin/bash

gostatic ()
{
    local dir=$1;
    local arg=$2;
    if [[ -z $dir ]]; then
        dir=$(pwd);
    fi;
    local name;
    name=$(basename "$dir");
    ( cd "$dir" || exit;
    export GOOS=linux;
    echo "Building static binary for $name in $dir";
    case $arg in
        "netgo")
            set -x;
            go build -a -tags 'netgo static_build' -installsuffix netgo -ldflags "-w" -o "$name" .
        ;;
        "cgo")
            set -x;
            CGO_ENABLED=1 go build -a -tags 'cgo static_build' -ldflags "-w -extldflags -static" -o "$name" .
        ;;
        *)
            set -x;
            CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-w" -o "$name" .
        ;;
    esac )
}

gostatic

