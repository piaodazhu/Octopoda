#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

if [ ! -f "brain" ]; then
  echo "brain not found."
  exit 1
fi
if [ ! -f "brain.yaml" ]; then
  echo "brain.yaml not found."
  exit 1
fi
if [ ! -f "brain.service" ]; then
  echo "brain.service not found."
  exit 1
fi

# binary
if [ ! -d "/usr/local/bin/octopoda" ]; then
  mkdir -p /usr/local/bin/octopoda
  echo "create folder /usr/local/bin/octopoda"
fi

# configuration
if [ ! -d "/etc/octopoda/brain" ]; then
  mkdir -p /etc/octopoda/brain
  echo "create folder /etc/octopoda/brain"
fi
cp brain /usr/local/bin/octopoda/
chmod +x /usr/local/bin/octopoda/brain
echo "install binary executable file --> /usr/local/bin/octopoda/"
cp brain.yaml /etc/octopoda/brain/
echo "install config brain.yaml --> /etc/octopoda/brain/"

# systemctl
cp brain.service /etc/systemd/system/
echo "create brain deamon"

systemctl enable brain
systemctl start brain
echo "start brain deamon"

echo ">> Setup Done"