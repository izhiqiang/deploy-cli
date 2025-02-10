#!/bin/bash

set -eux

#on:
#  workflow_dispatch:
#
#permissions:
#    contents: write
#    packages: write
#
#jobs:
#  release-linux-amd64:
#    name: release linux/amd64
#    runs-on: ubuntu-latest
#    steps:
#    - uses: actions/checkout@v4
#    - name: Set up Go
#      uses: actions/setup-go@v4
#      with:
#        go-version: '1.20'
#    - name: Build
#      run:  curl -sSL https://raw.githubusercontent.com/izhiqiang/deploy-cli/main/sh/actions-run.sh | bash
#      env:
#        DEPLOY_HOSTS: "ssh://root:123456@127.0.0.1"

TEMP="$(mktemp -d)"
trap 'rm -rf $TEMP' EXIT INT

URL="https://raw.githubusercontent.com/izhiqiang/deploy-cli/main/sh/install.sh"
wget --progress=dot:mega "$URL" -O "$TEMP/install.sh" || {
    echo "download installation install.sh failed"
    exit 1
}

source "$TEMP/install.sh"
DEPLOY_HOSTS="${DEPLOY_HOSTS:-""}"
deploy-cli run
