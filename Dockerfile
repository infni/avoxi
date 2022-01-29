FROM golang:1.17 as GOBASE

RUN apt-get update \
    && apt-get -y install unzip

WORKDIR /app

COPY . .

# Instal Protoc in the shell path
RUN PROTOC_VERSION="3.19.4" \
    && PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip" \
    && wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}" \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \
    && rm -f $PROTOC_ZIP

RUN go mod download

RUN go get github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor@v2.7.3
RUN go mod download google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN ./dev/build-service.sh \
    && go vet ./... \
    && result=`go test -timeout 2s -v ./test/... 2>&1` ; rc=$? \
    && echo "$result" | tee testResults.txt \
    && [ "$rc" -eq 0 ] \
    && go build -v -o "./bin/ipauthorize" -ldflags="-X main.version=$(git describe --always --long)" ./cmd/ipauthorize

RUN ./bin/ipauthorize -version > version && cat version

ARG PUB="/pub"
RUN mkdir -p $PUB \
    && cp ./bin/ipauthorize $PUB \
    && cp ./GeoLite2-Country.mmdb $PUB \
    && cp testResults.txt $PUB  \
    && cp version $PUB

# docker-in-docker
FROM gcr.io/distroless/base

WORKDIR /pub

COPY --from=GOBASE /pub .

CMD ["/pub/ipauthorize"]