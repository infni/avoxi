#! /bin/bash

# for package installation referece
# https://dev.to/calvinallen/installing-docker-and-docker-compose-in-wsl2ubuntu-on-windows-1c5b
# for service config referece
# https://dev.to/bowmanjd/install-docker-on-windows-wsl-without-docker-desktop-34m9

sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release

if [ ! /usr/share/keyrings/docker-archive-keyring.gpg ]; then
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
fi

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

sudo usermod -aG docker $USER

DOCKER_DIR=/mnt/wsl/shared-docker
sudo mkdir -pm o=,ug=rwx "$DOCKER_DIR"
sudo chgrp docker "$DOCKER_DIR"

sudo mkdir -p /etc/docker/

cat <<EOF >> ~/.bashrc

# start docker on login if it isn't already started
DOCKER_DISTRO="Ubuntu"
DOCKER_DIR=/mnt/wsl/shared-docker
DOCKER_SOCK="$DOCKER_DIR/docker.sock"
export DOCKER_HOST="unix://$DOCKER_SOCK"
if [ ! -S "$DOCKER_SOCK" ]; then
    mkdir -pm o=,ug=rwx "$DOCKER_DIR"
    chgrp docker "$DOCKER_DIR"
    /mnt/c/Windows/System32/wsl.exe -d $DOCKER_DISTRO sh -c "nohup sudo -b dockerd < /dev/null > $DOCKER_DIR/dockerd.log 2>&1"
fi
EOF

# this 'sudo' call wasn't working at last run, and I used nano to make the file.
sudo cat <<EOF > /etc/docker/daemon.json
{
  "hosts": ["unix:///mnt/wsl/shared-docker/docker.sock"],
  "iptables": false
}
EOF

# this 'sudo' call wasn't working at last run, and I used `sudo visudo`` to append to the file.
sudo cat <<EOF >> /etc/sudoers

# stop docker from asking for a password every time you log in
%docker ALL=(ALL)  NOPASSWD: /usr/bin/dockerd
EOF