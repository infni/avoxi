#!/bin/sh

PROTOC_VERSION="3.19.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip"
wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
rm -f $PROTOC_ZIP

go get github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor@v2.7.3
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc 
