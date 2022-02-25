#!/bin/sh

sudo apt-get update
sudo apt-get install -y protobuf-compiler zip unzip

GO_VERSION="1.17.7"
GO_FILENAME="go${GO_VERSION}.linux-amd64.tar.gz"
wget "https://dl.google.com/go/${GO_FILENAME}"
tar -xvf ${GO_FILENAME}
sudo mv go /usr/local
rm -f $GO_FILENAME

cat <<'EOF' >> ~/.bashrc

# Set up GO, and add GO to the path
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
EOF

source ~/.profile

PROTOC_VERSION="3.19.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip"
wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
#sudo chmod +x /usr/local/bin/protoc
rm -f $PROTOC_ZIP

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc 
