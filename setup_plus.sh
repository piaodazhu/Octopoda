#!/bin/bash
if [[ $# == 0 ]];
then
	echo "usage: bash setup_plus.sh <tentacle_name>, or bash setup_plus.sh tlist <namelist_file>"
	exit 0
fi

if [[ $# == 2 ]];
then
	for line in `cat $2`
	do
		bash setup_plus.sh $line
	done && \
	echo "DONE"
	exit 0
elif [[ $# > 2 ]];
then	
	echo "usage: bash setup_plus.sh <tentacle_name>, or bash setup_plus.sh tlist <namelist_file>"
	exit 1
fi

NAME=$1
INSTALL_ARCH=amd64
INSTALL_VERSION='1.8.4'
OUTPUT_DIR=./installers/

HTTPNS_URL=https://10.108.30.85:3455/release
CERTGEN_SCRIPT=./httpns/CertGen.sh
CACERT_DIR=./httpns/ca

# prepare scripts and keys
WORK_DIR=${OUTPUT_DIR}/installer_$NAME
mkdir -p ${WORK_DIR} && cp ${CERTGEN_SCRIPT} ${WORK_DIR} && cp ${CACERT_DIR}/ca.key ${CACERT_DIR}/ca.pem ${WORK_DIR}
cd ${WORK_DIR} && bash ./CertGen.sh client $NAME && mv ca.pem $NAME && rm -rf ca.key && rm -rf CertGen.sh || (rm -rf ${WORK_DIR} && exit 1)

# fetch installation package
wget --ca-certificate=$NAME/ca.pem --certificate=$NAME/client.pem --private-key=$NAME/client.key \
	${HTTPNS_URL}/tentacle_v${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)
wget --ca-certificate=$NAME/ca.pem --certificate=$NAME/client.pem --private-key=$NAME/client.key \
	${HTTPNS_URL}/pakma_v${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)

# write installation script
cat > install.sh << EOF
if [ "\$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi
tar -Jxvf tentacle_v${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz && \\
	cd tentacle_v${INSTALL_VERSION}_linux_${INSTALL_ARCH} && \\
	sed -i "s/pi0/$NAME/" tentacle.yaml && bash setup.sh && \\
	echo "setup tentacle done" || exit 1
cd ../ && mkdir -p /etc/octopoda/cert/ && cp $NAME/ca.pem $NAME/client.key $NAME/client.pem /etc/octopoda/cert/ && \\
	echo "copy cert and key done" || exit 1
mkdir -p /var/octopoda/tentacle/pakma/ && cp tentacle_v${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz /var/octopoda/tentacle/pakma/ && \\
	echo '{"StateType":2,"Version1":"${INSTALL_VERSION}","Version2":"${INSTALL_VERSION}","Version3":""}' > /var/octopoda/tentacle/pakma/pakma.json || exit 1
tar -Jxvf pakma_v${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz && \\
	cd pakma_v${INSTALL_VERSION}_linux_${INSTALL_ARCH} && bash setup.sh tentacle && \\
	echo "setup pakma_tentacle done" || exit 1
echo "DONE"
EOF

# pack the installation package 
cd ../
tar -Jcvf installer_$NAME.tar.xz installer_$NAME && rm -rf installer_$NAME && echo "pack installer_$NAME done" || exit 1
