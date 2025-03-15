#!/bin/bash

PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")
echo $PROJ_ROOT_DIR

OUTPUT_DIR=${PROJ_ROOT_DIR}/_output
echo $OUTPUT_DIR

VERSION_PACKAGE=github.com/onexstack_practice/fast_blog/pkg/version

if [[ -z "${VERSION}" ]];then 
	VERSION=$(git describe --tags --always --match='v*')
fi

GIT_TREE_STATE="dirty"

is_clean=$(shell git status --porcelain 2>/dev/null)

if [[ -z ${is_clean} ]];then
	GIT_TREE_STATE="clean"
fi

GIT_COMMIT=$(git rev-parse HEAD)
echo $GIT_COMMIT

GO_LDFALGS="-X ${VERSION_PACKAGE}.gitVersion=${VERSION} \
	-X ${VERSION_PACKAGE}.gitCommit=${GIT_COMMIT} \
	-X ${VERSION_PACKAGE}.gitTreeState=${GIT_TREE_STATE} \ 
	-X ${VERSION_PACKAGE}.buildDate=$(date -u +'%Y-%m-%d%H:%M:%SZ')"

go build -v -ldflags "${GO_LDFLAGS}" -o ${OUTPUT_DIR}/fb-apiserver -v cmd/fb-apiserver/main.go 
