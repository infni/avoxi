#!/bin/sh

sudo apt-get update
sudo apt-get install -y protobuf-compiler zip unzip gcc

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