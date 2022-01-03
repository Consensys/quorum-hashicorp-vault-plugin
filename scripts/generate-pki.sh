#!/bin/bash

# useful dirs
TEMP_DIR=/tmp

mkdir -p $DEST_CERT_PATH

ROOT_CRT=$TEMP_DIR/ca.pem
ROOT_KEY=$TEMP_DIR/ca-key.pem

# Gen root + intermediate
cfssl gencert -initca $CONF_DIR/root.json | cfssljson -bare $TEMP_DIR/ca
cfssl gencert -initca $CONF_DIR/intermediate.json | cfssljson -bare $TEMP_DIR/intermediate_ca
cfssl sign -ca $ROOT_CRT -ca-key $ROOT_KEY -config $CONF_DIR//cfssl.json -profile intermediate_ca $TEMP_DIR/intermediate_ca.csr | cfssljson -bare $TEMP_DIR/intermediate_ca

INTER_CRT=$TEMP_DIR/intermediate_ca.pem
INTER_KEY=$TEMP_DIR/intermediate_ca-key.pem

# Gen leaves from intermediate
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=server $CONF_DIR/vault.json | cfssljson -bare $TEMP_DIR/vault-server
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=client $CONF_DIR/vault.json | cfssljson -bare $TEMP_DIR/vault-client

# ca.crt is ca.pem >> intermediate.pem
cat $TEMP_DIR/ca.pem > $TEMP_DIR/ca.crt
cat $TEMP_DIR/intermediate_ca.pem >> $TEMP_DIR/ca.crt

# Verify certs
openssl verify -CAfile $TEMP_DIR/ca.crt $TEMP_DIR/vault-server.pem $TEMP_DIR/vault-client.pem

# Relocate
mv $TEMP_DIR/ca.crt $DEST_CERT_PATH/ca.crt
mv $TEMP_DIR/vault-server.pem $DEST_CERT_PATH/tls.crt
mv $TEMP_DIR/vault-server-key.pem $DEST_CERT_PATH/tls.key
mv $TEMP_DIR/vault-client.pem $DEST_CERT_PATH/client.crt
mv $TEMP_DIR/vault-client-key.pem $DEST_CERT_PATH/client.key

# cleanup
rm -rf $TEMP_DIR


