#!/bin/sh
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

if [ ! -d "/root/octopoda/brain/bin" ]; then
  mkdir -p /root/octopoda/brain/bin
  echo "create folder /root/octopoda/brain/bin"
fi

if [ ! -d "/etc/octopoda/brain" ]; then
  mkdir -p /etc/octopoda/brain
  echo "create folder /etc/octopoda/brain"
fi
cp brain /root/octopoda/brain/bin
chmod +x /root/octopoda/brain/bin/brain
echo "install binary executable file --> /root/octopoda/brain/bin/"
cp brain.yaml /etc/octopoda/brain/
echo "install config brain.yaml --> /etc/octopoda/brain/"

cp brain.service /etc/systemd/system/
echo "create brain deamon"

systemctl enable brain
systemctl start brain
echo "start brain deamon"

echo ">> Setup Done"