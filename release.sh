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

if [[ $# != 1 ]]; then
    printUsage "Usage: ./`basename $0` version"
    exit 1
fi

if ! [[ `echo "$1" | ag '^v\d+\.\d+\.\d+$'` ]]; then
    printError "Invalid version format: $1"
    exit 1
fi

if [[ `git tag | ag -Q "$1"` ]]; then
    printError "tag '$1' already exists"
    exit 1
fi

[[ -d release ]] && rm -rf release

function buildTools() {
    echo "Building for $1..."

    binaryFile="archivist"
    cd cli/archivist
    GOOS="$1" GOARCH=amd64 go build -o "$2/$binaryFile$3"
    [[ $? -ne 0 ]] && exit 1
    cd ../..

    binaryFile="easyjson"
    GOOS="$1" GOARCH=amd64 go build -o "$2/$binaryFile$3" github.com/edwingeng/easyjson-alt/easyjson
    [[ $? -ne 0 ]] && exit 1

    cp -r cli/script "$2/script"
    [[ $? -ne 0 ]] && exit 1

    git log --date=iso head~1..head > "$2/commit"
    [[ $? -ne 0 ]] && exit 1

    cd release
    zip -r -m "$4".zip `basename $2`
    [[ $? -ne 0 ]] && exit 1
    cd ..
    echo
}

goos=darwin
outputDir="`pwd`/release/archivist-$goos-$1"
buildTools "$goos" "$outputDir" "" "archivist-mac-$1"

goos=linux
outputDir="`pwd`/release/archivist-$goos-$1"
buildTools "$goos" "$outputDir" "" "archivist-$goos-$1"

goos=windows
outputDir="`pwd`/release/archivist-$goos-$1"
buildTools "$goos" "$outputDir" ".exe" "archivist-$goos-$1"

git tag "$1"
[[ $? -ne 0 ]] && exit 1
gitTag=`git tag | ag -Q "$1" | xargs echo "tag:"`
printImportantMessage "$gitTag"

echo ''
echo '====== done! ======'
