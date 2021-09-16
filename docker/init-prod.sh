# Store root token in a file so it can be shared with other services through volume
# Init Vault

echo "Initializing Vault: ${VAULT_ADDR}"

curl --request POST --data '{"secret_shares": 1, "secret_threshold": 1}' ${VAULT_ADDR}/v1/sys/init > init.json

if [ "$UNSEAL_KEY" = "null" ]; then
  echo "cannot retrieve unseal token"
  exit 1
fi

# Retrieve root token and unseal key
VAULT_TOKEN=$(cat init.json | jq .root_token | tr -d '"')
UNSEAL_KEY=$(cat init.json | jq .keys | jq .[0])
SHA256SUM=$(sha256sum -b ${PLUGIN_FILE} | cut -d' ' -f1)
rm init.json


# Unseal Vault
echo "unsealing vault..."
curl --request POST --data '{"key": '${UNSEAL_KEY}'}' ${VAULT_ADDR}/v1/sys/unseal

echo "registering Quorum Hashicorp Vault plugin..."
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data "{\"sha256\": \"${SHA256SUM}\", \"command\": \"quorum-hashicorp-vault-plugin\" }" \
  ${VAULT_ADDR}/v1/sys/plugins/catalog/secret/quorum-hashicorp-vault-plugin

echo "enabling Quorum Hashicorp Vault engine..."
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "plugin", "plugin_name": "quorum-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/quorum

echo "ROOT_TOKEN: $VAULT_TOKEN"
