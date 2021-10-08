[![Website](https://img.shields.io/website?label=documentation&url=https%3A%2F%2Fdocs.orchestrate.pegasys.tech%2F)](https://docs.orchestrate.pegasys.tech/)
[![Website](https://img.shields.io/website?label=website&url=https%3A%2F%2Fpegasys.tech%2Forchestrate%2F)](https://pegasys.tech/orchestrate/)
[![CircleCI](https://img.shields.io/circleci/build/gh/ConsenSys/quorum-hashicorp-vault-plugin?token=8a52ab8f0640f5bee56991cd30d808f735749dbf)](https://circleci.com/gh/PegaSysEng/quorum-hashicorp-vault-plugin)

![](/img/QuorumLogo_Blue.png)

# Quorum Hashicorp Vault plugin

The Quorum plugin enhances Hashicorp Vault Service with cryptographic operations under Vault engine, such as:
 - Create and import keys with the following supported eliptic curve and signing algorithm: ecdsa+sepc256k1 or eddsa+babyjubjub
 - Sign with every supported key pair. 
 - Create and import Ethereum wallets
 - Sign Ethereum transactions
 - Sign EEA private transaction
 - Sign Quorum Tessera private transaction
 - Create and import ZKP accounts
 - ZKP signing operation

## Development

### Pre-requirements
- Go >= 1.15
- Makefile
- docker-compose

### Running local version

Build plugin binary
```bash
$> make gobuild
```

To run our plugin in development mode you have to first build the plugin using:
```bash
$> make dev
```

### Testing

Now you have your Vault running on port `:8200`. Open  a new terminal to run the following command to
enable Orchestrate plugin:
```bash
$> curl --header "X-Vault-Token: DevVaultToken" --request POST \
  --data '{"type": "plugin", "plugin_name": "quorum-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/quorum
```

Now you already have your Vault running with Orchestrate plugin enable. The best way to understand the new
 integrate APIs is to use the `help` feature. To list a description of all the available endpoints you can run:
```bash
$> curl -H "X-Vault-Token: DevVaultToken" http://127.0.0.1:8200/v1/quorum?help=1
```

alternatively you can list only `ethereum` endpoints by using:
```bash
$> curl -H "X-Vault-Token: DevVaultToken" http://127.0.0.1:8200/v1/quorum/ethereum/accounts?help=1
```

## Running using latest version

Running Quorum Hashicorp Vault Plugin plugin:
```bash
$> docker-compose -f docker-compose.yml up --build vault
```

## Contributing
[How to Contribute](CONTRIBUTING.md)

## License

Quorum Hashicorp Vault plugin is licensed under the BSL 1.1.

Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [quorum-key-manager@consensys.net](mailto:quorum-key-manager@consensys.net) if you need to purchase a license for a production use-case.  

