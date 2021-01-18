echo "enabling orchestrate secret engine..."
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "plugin", "plugin_name": "orchestrate-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/orchestrate
