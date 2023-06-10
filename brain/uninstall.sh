#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

systemctl disable brain
systemctl stop brain
echo "stop brain deamon"

if [ -f "/etc/systemd/system/brain.service" ]; then
  rm /etc/systemd/system/brain.service
  echo "remove file /etc/systemd/system/brain.service"
fi

if [ -f "/usr/local/bin/octopoda/brain" ]; then
  rm -rf /usr/local/bin/octopoda/brain
  echo "remove executable /usr/local/bin/octopoda/brain"
fi

if [ -d "/etc/octopoda/brain/brain.yaml" ]; then
  rm -rf /etc/octopoda/brain/brain.yaml
  echo "remove configuration /etc/octopoda/brain/brain.yaml"
fi

echo ">> Uninstall Done"