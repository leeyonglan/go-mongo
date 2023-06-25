#!/bin/bash -e

DIR=$(readlink -f "$0") && DIR=$(dirname "$DIR") && cd "$DIR"

DATE=$(TZ='Asia/Shanghai' date '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(go version)

source ./common.sh

LDFLAGS="-X '${BUILD_PACKAGE}.BuildGoVersion=${GO_VERSION}' \
	-X '${BUILD_PACKAGE}.BuildTime=${DATE}' \
	-X '${BUILD_PACKAGE}.BuildType=${TYPE}'"

cd ..

if [ -d "vendor" ]; then
	VENDOR="-mod=vendor"
fi

if [ "${TYPE}" = "dev" ]; then
    echo "dev"
    CGO_ENABLED=0 go build $VENDOR \
	-ldflags "$LDFLAGS" \
	-o "$EXE" \
    2>&1
else
    echo "pro"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $VENDOR \
	-ldflags "$LDFLAGS" \
	-o "$EXE" \
    2>&1
fi