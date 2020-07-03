#!/bin/bash
#

set -eu
date

apt update
sudo apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common

curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo add-apt-repository   "deb [arch=arm64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"

apt-get update
apt-get install docker-ce docker-ce-cli containerd.io
if [[ $? != 0 ]]; then
    echo "error install"
    exit 255
fi
echo "Done"