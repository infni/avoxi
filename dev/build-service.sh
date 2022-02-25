#!/bin/sh

protoc -I. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/IpAuthorize.proto
protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:. api/IpAuthorize.proto

go build -v -o ./bin/ipauthorize -ldflags="-X main.version=$(git describe --always --long)" ./cmd/ipauthorize
