#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

if [ ! -f "httpns" ]; then
  echo "httpns not found."
  exit 1
fi
if [ ! -f "httpns.yaml" ]; then
  echo "httpns.yaml not found."
  exit 1
fi
if [ ! -f "httpns.service" ]; then
  echo "httpns.service not found."
  exit 1
fi

# binary
if [ ! -d "/usr/local/bin/octopoda" ]; then
  mkdir -p /usr/local/bin/octopoda
  echo "create folder /usr/local/bin/octopoda"
fi

# configuration
if [ ! -d "/etc/octopoda/httpns" ]; then
  mkdir -p /etc/octopoda/httpns
  echo "create folder /etc/octopoda/httpns"
fi
cp httpns /usr/local/bin/octopoda/
chmod +x /usr/local/bin/octopoda/httpns
echo "install binary executable file --> /usr/local/bin/octopoda/"
cp httpns.yaml /etc/octopoda/httpns/
echo "install config httpns.yaml --> /etc/octopoda/httpns/"

# systemctl
cp httpns.service /etc/systemd/system/
echo "create httpns deamon"

systemctl enable httpns
systemctl start httpns
echo "start httpns deamon"

echo ">> Setup Done"