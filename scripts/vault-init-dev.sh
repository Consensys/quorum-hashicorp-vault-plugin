VAULT_ADDR=${VAULT_ADDR-localhost:8200}
PLUGIN_MOUNT_PATH=${PLUGIN_MOUNT_PATH-quorum}
PLUGIN_PATH=${PLUGIN_PATH-/vault/plugins}
VAULT_DEV_ROOT_TOKEN_ID=${VAULT_DEV_ROOT_TOKEN_ID-DevVaultToken}

if [ "${PLUGIN_PATH}" != "/vault/plugins" ]; then
  mkdir -p ${PLUGIN_PATH}
  echo "[PLUGIN] Copying plugin to expected folder"
  cp $PLUGIN_FILE "${PLUGIN_PATH}/quorum-hashicorp-vault-plugin"
fi 

echo "[PLUGIN] Enabling Quorum Hashicorp Plugin engine..."
curl --header "X-Vault-Token: ${VAULT_DEV_ROOT_TOKEN_ID}" --request POST \
  --data '{"type": "plugin", "plugin_name": "quorum-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/${PLUGIN_MOUNT_PATH}
  
