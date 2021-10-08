echo "enabling Quorum Hashicorp Plugin engine..."
curl --header "X-Vault-Token: ${VAULT_DEV_ROOT_TOKEN_ID}" --request POST \
  --data '{"type": "plugin", "plugin_name": "quorum-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR-http://localhost:8200}/v1/sys/mounts/quorum
