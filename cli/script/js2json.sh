#!/usr/bin/env bash

##################################################
# Owned by watcher. DON'T change me.
##################################################

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

if [[ $# -ne 1 ]]; then
    printUsage "Usage: ./`basename $0` [RootDir]"
    exit 1
fi

if [[ "$1" != /* ]]; then
    printError "RootDir must be an absolute directory."
    exit 1
fi

if [[ "$WATCHER_CONF_GROUP" == '' ]]; then
    WATCHER_CONF_GROUP=local
fi
if [[ ! -d "$1/$WATCHER_CONF_GROUP" ]]; then
    printError "$1/$WATCHER_CONF_GROUP does not exist."
    exit 1
fi

printImportantMessage "[#] Converting .js to its corresponding .json..."
echo
echo "WATCHER_CONF_GROUP: $WATCHER_CONF_GROUP"
echo "WATCHER_CONF_SUBGROUP: $WATCHER_CONF_SUBGROUP"
echo

rm -rf "$1/.runtime"
if [[ "$WATCHER_CONF_SUBGROUP" == "" ]]; then
    mkdir -p "$1/.runtime/$WATCHER_CONF_GROUP"
else
    mkdir -p "$1/.runtime/$WATCHER_CONF_GROUP/$WATCHER_CONF_SUBGROUP"
fi

for f in `find "$1" -maxdepth 1 -name '*.js'`; do
    echo "Converting $f..."
    g=`basename $f .js`
    h="$1/.runtime/$g.json"
    node js2json.js "$f" > "$h"
    [[ $? -ne 0 ]] && exit 1
done

for f in `find "$1/$WATCHER_CONF_GROUP" -maxdepth 1 -name '*.js'`; do
    echo "Converting $f..."
    g=`basename $f .js`
    h="$1/.runtime/$WATCHER_CONF_GROUP/$g.json"
    node js2json.js "$f" > "$h"
    [[ $? -ne 0 ]] && exit 1
done

if [[ "$WATCHER_CONF_SUBGROUP" != "" ]]; then
    for f in `find "$1/$WATCHER_CONF_GROUP/$WATCHER_CONF_SUBGROUP" -maxdepth 1 -name '*.js'`; do
        echo "Converting $f..."
        g=`basename $f .js`
        h="$1/.runtime/$WATCHER_CONF_GROUP/$WATCHER_CONF_SUBGROUP/$g.json"
        node js2json.js "$f" > "$h"
        [[ $? -ne 0 ]] && exit 1
    done
fi

:
