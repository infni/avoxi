#!/bin/sh -v

wget https://dot.net/v1/dotnet-install.sh -O dotnet-install.sh
chmod +x ./dotnet-install.sh
./dotnet-install.sh --channel 8.0
dotnet --version

cat <<'EOF' >> ~/.bashrc

# Set up dotnet on the path
export PATH=/home/infni7/.dotnet:$PATH
EOF

source ~/.profile
