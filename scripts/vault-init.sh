# Store root token in a file so it can be shared with other services through volume
# Init Vault

VAULT_ADDR=${VAULT_ADDR-localhost:8200}
PLUGIN_PATH=${PLUGIN_PATH-/vault/plugins}
PLUGIN_MOUNT_PATH=${PLUGIN_MOUNT_PATH-quorum}
ROOT_TOKEN_PATH=${ROOT_TOKEN_PATH-/vault/.root}
PLUGIN_FILE=/vault/plugins/quorum-hashicorp-vault-plugin

VAULT_SSL_PARAMS=""
if [ -n "$VAULT_CACERT" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --cacert $VAULT_CACERT"
fi  

if [ -n "$VAULT_CLIENT_CERT" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --cert $VAULT_CLIENT_CERT"
fi     

if [ -n "$VAULT_CLIENT_KEY" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --key $VAULT_CLIENT_KEY"
fi     

echo "[PLUGIN] Initializing Vault: ${VAULT_ADDR}"

curl -s --request POST ${VAULT_SSL_PARAMS} \
  --data '{"secret_shares": 1, "secret_threshold": 1}' ${VAULT_ADDR}/v1/sys/init > response.json

ROOT_TOKEN=$(cat response.json | jq .root_token | tr -d '"')
UNSEAL_KEY=$(cat response.json | jq .keys | jq .[0])
ERRORS=$(cat response.json | jq .errors | jq .[0])
rm response.json

if [ "$UNSEAL_KEY" = "null" ]; then
  echo "[PLUGIN] cannot retrieve unseal key: $ERRORS"
  exit 1
fi

# Unseal Vault
echo "[PLUGIN] Unsealing vault..."
curl -s --request POST ${VAULT_SSL_PARAMS} \
  --data '{"key": '${UNSEAL_KEY}'}' ${VAULT_ADDR}/v1/sys/unseal

if [ "${PLUGIN_PATH}" != "/vault/plugins" ]; then
  mkdir -p ${PLUGIN_PATH}
  echo "[PLUGIN] Copying plugin to expected folder"
  cp $PLUGIN_FILE "${PLUGIN_PATH}/quorum-hashicorp-vault-plugin"
fi 

echo "[PLUGIN] Registering Quorum Hashicorp Vault plugin..."
SHA256SUM=$(sha256sum -b ${PLUGIN_FILE} | cut -d' ' -f1)
curl -s --header "X-Vault-Token: ${ROOT_TOKEN}" --request POST ${VAULT_SSL_PARAMS} \
  --data "{\"sha256\": \"${SHA256SUM}\", \"command\": \"quorum-hashicorp-vault-plugin\" }" \
  ${VAULT_ADDR}/v1/sys/plugins/catalog/secret/quorum-hashicorp-vault-plugin

echo "[PLUGIN] Enabling Quorum Hashicorp Vault engine..."
curl -s --header "X-Vault-Token: ${ROOT_TOKEN}" --request POST ${VAULT_SSL_PARAMS} \
  --data '{"type": "plugin", "plugin_name": "quorum-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/${PLUGIN_MOUNT_PATH}

if [ -n "$VAULT_KVV2_MOUNT_PATH" ]; then
# Enable kv-v2 secret engine
echo "[PLUGIN] Enabling kv-v2 Hashicorp Vault engine..."
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST ${VAULT_SSL_PARAMS}\
     --data '{"type": "kv-v2", "config": {"force_no_cache": true} }' \
     ${VAULT_ADDR}/v1/sys/mounts/${VAULT_KVV2_MOUNT_PATH}
fi

if [ -n "$ROOT_TOKEN" ]; then 
  echo "[PLUGIN] Root token saved in ${ROOT_TOKEN_PATH}"
  echo "$ROOT_TOKEN" > ${ROOT_TOKEN_PATH}
fi

exit 0
