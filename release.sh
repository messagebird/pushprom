#!/usr/bin/env bash
# a hack to generate releases like other prometheus projects
# use like this:
#       VERSION=1.0.1 ./release.sh
set -e

# github user and repo
USER=messagebird
REPO=pushprom
BIN_PACKAGE=github.com/${USER}/${REPO}
BIN_DIR=bin

if [ -z "$VERSION" ]; then
    >&2 echo "missing VERSION=X.X.X"
    exit 1
fi

function build_for {
    BIN_NAME="bin/${REPO}-${VERSION}.$1-$2"
    GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -ldflags "-s -w" -o ${BIN_NAME} ${BIN_PACKAGE}
    shasum -a 256 ${BIN_NAME} > ${BIN_NAME}.sha256
}

rm -rf ${BIN_DIR}
mkdir -p ${BIN_DIR}

build_for linux amd64
build_for darwin amd64

DOCKER_BUILDKIT=1 docker build -t ${USER}/${REPO}:${VERSION} .

git tag -a $VERSION -m "version $VERSION"

github-release release \
    --user $USER \
    --repo $REPO \
    --tag $VERSION \
    --name $VERSION \
    --description "version $VERSION"

for release_file in $(ls ${BIN_DIR}); do
    github-release upload \
        --user $USER \
        --repo $REPO \
        --tag $VERSION \
        --name "${release_file}" \
        --file "${BIN_DIR}/${release_file}"
done

docker push ${USER}/${REPO}:${VERSION}
docker tag ${USER}/${REPO}:${VERSION} ${USER}/${REPO}:latest
docker push ${USER}/${REPO}:latest
