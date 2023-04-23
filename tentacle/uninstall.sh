#!/bin/sh
systemctl disable tentacle
systemctl stop tentacle
echo "stop tentacle deamon"

if [ -f "/etc/systemd/system/tentacle.service" ]; then
  rm /etc/systemd/system/tentacle.service
  echo "remove file /etc/systemd/system/tentacle.service"
fi

if [ -d "/root/octopoda/tentacle/bin" ]; then
  rm -rf /root/octopoda/tentacle/bin
  echo "remove folder /root/octopoda/tentacle/bin"
fi

if [ -d "/etc/octopoda" ]; then
  rm -rf /etc/octopoda
  echo "remove folder /etc/octopoda"
fi

echo ">> Uninstall Done"