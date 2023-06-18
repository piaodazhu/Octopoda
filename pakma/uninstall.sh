#!/bin/sh
APP=""
if [[ $# != 1 ]];
then
    echo "You must specific the app to be managed by packma! (Usage: bash setup.sh [tentacle|brain])"
    exit 1
elif [[ $1 == "tentacle" ]];
then
    APP=tentacle
elif [[ $1 == "brain" ]];
then
    APP=brain
else
    echo "now packma only support tentacle and brain! (Usage: bash setup.sh [tentacle|brain])"
fi

if [ "$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi

systemctl disable pakma_$APP
systemctl stop pakma_$APP
echo "stop pakma_$APP deamon"

if [ -f "/etc/systemd/system/pakma_$APP.service" ]; then
  rm /etc/systemd/system/pakma_$APP.service
  echo "remove service /etc/systemd/system/pakma_$APP.service"
fi

if [ -f "/usr/local/bin/octopoda/pakma_$APP" ]; then
  rm -rf /usr/local/bin/octopoda/pakma_$APP
  echo "remove executable /usr/local/bin/octopoda/pakma_$APP"
fi

echo ">> Uninstall Done"