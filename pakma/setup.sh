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

if [ ! -f "pakma" ]; then
  echo "pakma not found."
  exit 1
fi
if [ ! -f "pakma_$APP.service" ]; then
  echo "pakma_$APP.service not found."
  exit 1
fi

# binary
if [ ! -d "/usr/local/bin/octopoda" ]; then
  mkdir -p /usr/local/bin/octopoda
  echo "create folder /usr/local/bin/octopoda"
fi

# configuration
if [ ! -d "/etc/octopoda/$APP" ]; then
  mkdir -p /etc/octopoda/$APP
  echo "warning: $APP may haven't been installed"
fi

cp pakma /usr/local/bin/octopoda/pakma_$APP
chmod +x /usr/local/bin/octopoda/pakma_$APP
echo "install binary executable file pakma_$APP --> /usr/local/bin/octopoda/"

# systemctl
cp pakma_$APP.service /etc/systemd/system/
echo "create pakma_$APP deamon"

systemctl enable pakma_$APP
systemctl start pakma_$APP
echo "start pakma_$APP deamon"

echo ">> Setup Done"