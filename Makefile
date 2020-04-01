BINARY_NAME=http_guldan

VERSION=$(shell git describe --tags)
BUILD_TIME=`date +%FT%T%z`
GIT_HASH=`git rev-parse HEAD`

LDFLAGS="-X main.AppVersion=${VERSION} -X main.AppBuildTime=${BUILD_TIME} -X main.AppGitHash=${GIT_HASH}"

build: ## Build the binary
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} main.go

DEB_PACKAGE_NAME=$(BINARY_NAME)
DEB_PACKAGE_PREFIX=/usr/local/http_guldan
DEB_PACKAGE_DESCRIPTION="http_guldan download service"

deb:   ## Build deb package
	exec ./build-deb.sh $(DEB_PACKAGE_NAME) $(DEB_PACKAGE_PREFIX) $(DEB_PACKAGE_DESCRIPTION)

clean: ## Clean this build
	rm -rf ${BINARY_NAME}
	rm -rf build
	rm -rf *.deb

help:  ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: all
all: build

