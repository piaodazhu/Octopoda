#!/bin/sh
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

if [ ! -d "/root/octopoda/tentacle/bin" ]; then
  mkdir -p /root/octopoda/tentacle/bin
  echo "create folder /root/octopoda/tentacle/bin"
fi

if [ ! -d "/etc/octopoda" ]; then
  mkdir -p /etc/octopoda
  echo "create folder /etc/octopoda"
fi
cp tentacle /root/octopoda/tentacle/bin
chmod +x /root/octopoda/tentacle/bin/tentacle
echo "install binary executable file --> /root/octopoda/tentacle/bin/"
cp tentacle.yaml /etc/octopoda/
echo "install config tentacle.yaml --> /etc/octopoda/"

cp tentacle.service /etc/systemd/system/
echo "create tentacle deamon"

systemctl enable tentacle
systemctl start tentacle
echo "start tentacle deamon"

echo ">> Setup Done"