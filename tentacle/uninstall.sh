#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

systemctl disable tentacle
systemctl stop tentacle
echo "stop tentacle deamon"

if [ -f "/etc/systemd/system/tentacle.service" ]; then
  rm /etc/systemd/system/tentacle.service
  echo "remove service /etc/systemd/system/tentacle.service"
fi

if [ -f "/usr/local/bin/octopoda/tentacle" ]; then
  rm -rf /usr/local/bin/octopoda/tentacle
  echo "remove executable /usr/local/bin/octopoda/tentacle"
fi

if [ -f "/etc/octopoda/tentacle/tentacle.yaml" ]; then
  rm -rf /etc/octopoda/tentacle/tentacle.yaml
  echo "remove configuration /etc/octopoda/tentacle/tentacle.yaml"
fi

echo ">> Uninstall Done"