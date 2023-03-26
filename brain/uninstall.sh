#!/bin/sh
systemctl disable brain
systemctl stop brain
echo "stop brain deamon"

if [ ! -f "/etc/systemd/system/brain.service" ]; then
  rm /etc/systemd/system/brain.service
  echo "remove file /etc/systemd/system/brain.service"
fi

if [ ! -d "/root/octopoda" ]; then
  rm -rf /root/octopoda
  echo "remove folder /root/octopoda"
fi

if [ ! -d "/etc/octopoda" ]; then
  rm -rf /etc/octopoda
  echo "remove folder /etc/octopoda"
fi

echo ">> Uninstall Done"