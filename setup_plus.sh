#!/bin/bash
NAME=renjiang
INSTALL_VERSION=v1.5.2
INSTALL_ARCH=arm
HTTPNS_URL=https://10.108.30.85:3455/release
CERTGEN_SCRIPT=./httpNameServer/CertGen.sh
CACERT_DIR=./httpNameServer/ca
OUTPUT_DIR=./installers/


# prepare scripts and keys
WORK_DIR=${OUTPUT_DIR}/installer_$NAME
mkdir -p ${WORK_DIR} && cp ${CERTGEN_SCRIPT} ${WORK_DIR} && cp ${CACERT_DIR}/ca.key ${CACERT_DIR}/ca.pem ${WORK_DIR}
cd ${WORK_DIR} && bash ./CertGen.sh client $NAME && mv ca.pem $NAME && rm -rf ca.key && rm -rf CertGen.sh || (rm -rf ${WORK_DIR} && exit 1)

# fetch installation package
wget --ca-certificate=$NAME/ca.pem --certificate=$NAME/client.pem --private-key=$NAME/client.key \
	${HTTPNS_URL}/tentacle_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)
wget --ca-certificate=$NAME/ca.pem --certificate=$NAME/client.pem --private-key=$NAME/client.key \
	${HTTPNS_URL}/pakma_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)

# wget --ca-certificate /etc/octopoda/cert/ca.pem --certificate /etc/octopoda/cert/client.pem --certificate-type PEM --private-key=/etc/octopoda/cert/client.key ${HTTPNS_URL}/tentacle_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)
# wget --ca-certificate /etc/octopoda/cert/ca.pem --certificate /etc/octopoda/cert/client.pem --certificate-type PEM --private-key=/etc/octopoda/cert/client.key ${HTTPNS_URL}/pakma_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz || (rm -rf ${WORK_DIR} && exit 1)

# write installation script
cat > install.sh << EOF
if [ "\$(id -u)" != "0" ]; then
   echo "You must run this script as root" 1>&2
   exit 1
fi
tar -Jxvf tentacle_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz && \\
	cd tentacle_${INSTALL_VERSION}_linux_${INSTALL_ARCH} && \\
	sed -i "s/pi0/$NAME/" tentacle.yaml && bash setup.sh && \\
	echo "setup tentacle done" || exit 1
cd ../ && cp $NAME/ca.pem $NAME/client.key $NAME/client.pem /etc/octopoda/tentacle/cert/ && echo "copy cert and key done" || exit 1
tar -Jxvf pakma_${INSTALL_VERSION}_linux_${INSTALL_ARCH}.tar.xz && \\
	cd pakma_${INSTALL_VERSION}_linux_${INSTALL_ARCH} && bash setup.sh tentacle \\
	echo "setup pakma_tentacle done" || exit 1
echo "DONE"
EOF

# pack the installation package 
cd ../
tar -Jcvf installer_$NAME.tar.xz installer_$NAME && rm -rf installer_$NAME && echo "pack installer_$NAME done" || exit 1
