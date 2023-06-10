#!/bin/sh
if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

if [ ! -f "tentacle" ]; then
  echo "tentacle not found."
  exit 1
fi
if [ ! -f "tentacle.yaml" ]; then
  echo "tentacle.yaml not found."
  exit 1
fi
if [ ! -f "tentacle.service" ]; then
  echo "tentacle.service not found."
  exit 1
fi

# binary
if [ ! -d "/usr/local/bin/octopoda" ]; then
  mkdir -p /usr/local/bin/octopoda
  echo "create folder /usr/local/bin/octopoda"
fi

# configuration
if [ ! -d "/etc/octopoda/tentacle" ]; then
  mkdir -p /etc/octopoda/tentacle
  echo "create folder /etc/octopoda/tentacle"
fi
cp tentacle /usr/local/bin/octopoda/
chmod +x /usr/local/bin/octopoda/tentacle
echo "install binary executable file --> /usr/local/bin/octopoda/"
cp tentacle.yaml /etc/octopoda/tentacle/
echo "install config tentacle.yaml --> /etc/octopoda/tentacle/"

# systemctl
cp tentacle.service /etc/systemd/system/
echo "create tentacle deamon"

systemctl enable tentacle
systemctl start tentacle
echo "start tentacle deamon"

echo ">> Setup Done"