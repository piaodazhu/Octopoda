#!/bin/sh
if [ ! -d "/root/octopoda/workspace" ]; then
  mkdir -p /root/octopoda/workspace
  echo "create folder /root/octopoda/workspace"
fi

if [ ! -d "/root/octopoda/log" ]; then
  mkdir -p /root/octopoda/log
  echo "create folder /root/octopoda/log"
fi

if [ ! -d "/etc/octopoda" ]; then
  mkdir -p /etc/octopoda
  echo "create folder /etc/octopoda"
fi
cp brain /root/octopoda/
echo "install binary executable file --> /root/octopoda/"
cp brain.yaml /etc/octopoda/
echo "install config brain.yaml --> /etc/octopoda/"

cp brain.service /etc/systemd/system/
echo "create brain deamon"

systemctl enable brain
systemctl start brain
echo "start brain deamon"

echo ">> Setup Done"