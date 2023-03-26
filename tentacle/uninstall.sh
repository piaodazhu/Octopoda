#!/bin/sh
systemctl disable tentacle
systemctl stop tentacle
echo "stop tentacle deamon"

if [ ! -f "/etc/systemd/system/tentacle.service" ]; then
  rm /etc/systemd/system/tentacle.service
  echo "remove file /etc/systemd/system/tentacle.service"
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