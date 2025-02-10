#!/bin/bash

# curl -sSL https://raw.githubusercontent.com/izhiqiang/deploy-cli/main/sh/install.sh | bash

DEPLOY_CLI_VERSION="${DEPLOY_CLI_VERSION:-""}"

if ! command -v curl &> /dev/null || ! command -v jq &> /dev/null; then
  echo "Error: 'curl' and 'jq' are required but not installed."
  exit 1
fi

if [ -z "$DEPLOY_CLI_VERSION" ]; then
  echo "DEPLOY_CLI_VERSION is empty. Fetching the latest version from API..."
  LATEST_VERSION=$(curl --silent "https://api.github.com/repos/izhiqiang/deploy-cli/tags" | jq -r '.[0].name')
  if [ -n "$LATEST_VERSION" ]; then
    DEPLOY_CLI_VERSION="$LATEST_VERSION"
    echo "Fetched latest version: $DEPLOY_CLI_VERSION"
  else
    echo "Failed to fetch the latest version. Please check the API."
    exit 1
  fi
else
  echo "Using specified version: $DEPLOY_CLI_VERSION"
fi

os=$(uname -s | awk '{print tolower($0)}')
case $(uname -m) in
    "x86_64" | "amd64")
        ARCH="amd64" ;;
    "i386" | "i486" | "i586")
        ARCH="i386" ;;
    "aarch64" | "arm64")
        ARCH="arm64" ;;
    "armv6l" | "armv7l")
        ARCH="arm" ;;
    "s390x")
        ARCH="s390x" ;;
    *)
        ARCH=$(uname -m)
esac

sudo_command="sudo"
if ! command -v sudo &> /dev/null ; then
  sudo_command=""
fi

TEMP="$(mktemp -d)"
trap 'rm -rf $TEMP' EXIT INT

URL="https://github.com/izhiqiang/deploy-cli/releases/download/${DEPLOY_CLI_VERSION}/deploy-cli_${os}_${ARCH}.tar.gz"
echo "Downloading deploy-cli binary from $URL..."
wget --progress=dot:mega "$URL" -O "$TEMP/deploy-cli.tar.gz" || {
    echo "Failed to download deploy-cli binary. Please check your architecture or version."
    exit 1
}

echo "Extracting deploy-cli binary..."
BIN_PATH=/usr/local/deploy
$sudo_command rm -rf $BIN_PATH
$sudo_command mkdir -p $BIN_PATH
$sudo_command tar -zxf "$TEMP/deploy-cli.tar.gz" -C $BIN_PATH || {
    echo "Failed to extract deploy-cli binary. Check the tarball integrity."
    exit 1
}

echo "Command installation to $BIN_PATH"

export PATH=$BIN_PATH:$PATH
echo "Installation complete. Verifying..."
deploy-cli --help || echo "deploy-cli binary is not functioning correctly. Check installation logs."
