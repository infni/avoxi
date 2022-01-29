#!/bin/sh

PROTOC_VERSION="3.19.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip" \
    && wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}" \
    && sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \
    && rm -f $PROTOC_ZIP