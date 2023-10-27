#/bin/bash
TYPE=$1
CA_KEY=ca.key
CA_CRT=ca.pem
CN=""
FNAME=""
EXTARG=""

if [[ ($TYPE == "ca") && $# == 2 ]];
then
	FNAME=ca
	CN=$2

	openssl genrsa -out $FNAME.key 2048 && \
	openssl req -new -x509 -sha256 -days 3650 -key $FNAME.key -subj "/C=CN/ST=BJ/L=BJ/O=BIT/CN=$CN" -out $FNAME.pem && \
	echo "DONE"
	rm -rf $FNAME.srl
	exit 0
fi

if [[ ($TYPE == "clientlist" || $TYPE == "clist") && $# == 2 ]];
then
	for line in `cat $2`
	do
		bash CertGen.sh cli $line
	done && \
	echo "DONE"
	exit 0
fi

if [[ ($TYPE == "server" || $TYPE == "svr") && $# == 3 ]];
then
	FNAME=server
	EXTARG=subjectAltName=DNS:localhost,IP:127.0.0.1,IP:$3
elif [[ ($TYPE == "client" || $TYPE == "cli") && $# == 2 ]];
then
	FNAME=client
else
	echo "usage: ./CertGen.sh [ca|client|server] {CommonName} {IP_if_server}"
	exit 0
fi
# echo $EXTARG
# exit 0
CN=$2
OUTPUT=$CN

mkdir -p $OUTPUT && \
echo ">> generating key pair and cert signing request..." && \
openssl req -newkey rsa:2048 -nodes -sha256 -keyout $OUTPUT/$FNAME.key \
	-subj "/C=CN/ST=BJ/L=BJ/O=BIT/CN=$CN" -out $OUTPUT/$FNAME.csr &&\
echo ">> CA is signing cert..." && \
openssl x509 -req -extfile <(printf "$EXTARG") -sha256 -days 3650 -in $OUTPUT/$FNAME.csr \
	-CA $CA_CRT -CAkey $CA_KEY -CAcreateserial -out $OUTPUT/$FNAME.pem && \
rm $OUTPUT/$FNAME.csr && \
cp $CA_CRT $OUTPUT/ && \
echo "DONE"