#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

systemctl disable httpns
systemctl stop httpns
echo "stop httpns deamon"

if [ -f "/etc/systemd/system/httpns.service" ]; then
  rm /etc/systemd/system/httpns.service
  echo "remove service /etc/systemd/system/httpns.service"
fi

if [ -f "/usr/local/bin/octopoda/httpns" ]; then
  rm -rf /usr/local/bin/octopoda/httpns
  echo "remove executable /usr/local/bin/octopoda/httpns"
fi

if [ -f "/etc/octopoda/httpns/httpns.yaml" ]; then
  rm -rf /etc/octopoda/httpns/httpns.yaml
  echo "remove configuration /etc/octopoda/httpns/httpns.yaml"
fi

echo ">> Uninstall Done"