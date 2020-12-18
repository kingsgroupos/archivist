#!/usr/bin/env bash

[[ "$TRACE" ]] && set -x
pushd `dirname "$0"` > /dev/null
trap __EXIT EXIT

colorful=false
tput setaf 7 > /dev/null 2>&1
if [[ $? -eq 0 ]]; then
    colorful=true
fi

function __EXIT() {
    popd > /dev/null
}

function printError() {
    $colorful && tput setaf 1
    >&2 echo "Error: $@"
    $colorful && tput setaf 7
}

function printImportantMessage() {
    $colorful && tput setaf 3
    >&2 echo "$@"
    $colorful && tput setaf 7
}

function printUsage() {
    $colorful && tput setaf 3
    >&2 echo "$@"
    $colorful && tput setaf 7
}

NEWPATH=`../../../../privatePath.sh`
[[ $? -ne 0 ]] && exit 1
PATH="$NEWPATH"

WATCHER_CONF_GROUP=develop WATCHER_CONF_SUBGROUP=dolores ../../../../script/js2json.sh "`pwd`/json"
echo

go run ../../../cli/archivist/archivist.go generate --outputDir=conf --pkg=conf --x-easyjson "$@" 'json/*.json' 'json/*.js'
