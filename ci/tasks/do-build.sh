#!/bin/sh
basedir=`pwd`/gopath-tested/src/github.com/ecsteam/docker-usage
build_dir=`pwd`/build-output/build
version_file=`pwd`/version/number

mkdir ${build_dir} > /dev/null 2>&1

set -e
set -x

export GOPATH=`pwd`/gopath-tested

# Run tests
cd ${basedir}
for os in linux windows darwin; do
    suffix=${os}
    if [ "windows" = "${os}" ]; then
        suffix="windows.exe"
    elif [ "darwin" = "${os}" ]; then
        suffix="macosx"
    fi

    GOOS=${os} GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X github.com/ecsteam/docker-usage/command.version=`cat ${version_file}`" -o ${build_dir}/docker-usage-${suffix}
done
