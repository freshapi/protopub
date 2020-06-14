#!/bin/bash
VERSION=$(git describe --tags 2> /dev/null)

if [[ -z "$VERSION" ]]; then
  VERSION=$(git rev-parse --short HEAD)-dev
fi

export BUILD_VERSION=$VERSION

echo BUILD_VERSION $BUILD_VERSION
