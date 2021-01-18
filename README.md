[![Website](https://img.shields.io/website?label=documentation&url=https%3A%2F%2Fdocs.orchestrate.pegasys.tech%2F)](https://docs.orchestrate.pegasys.tech/)
[![Website](https://img.shields.io/website?label=website&url=https%3A%2F%2Fpegasys.tech%2Forchestrate%2F)](https://pegasys.tech/orchestrate/)
[![CircleCI](https://img.shields.io/circleci/build/gh/ConsenSys/orchestrate-hashicorp-vault-plugin?token=8a52ab8f0640f5bee56991cd30d808f735749dbf)](https://circleci.com/gh/PegaSysEng/orchestrate-hashicorp-vault-plugin)

![](/img/codefi_orchestrate_logo.png)

# Orchestrate Hashicorp Vault plugin

The Orchestrate library provides a secure Hashicorp Vault plugin to store your keys and perform
 cryptographic operations such as:
 - Ethereum transactions
 - Besu-Orion private transaction
 - Quorum-Tessera private transaction
 - ZKP signing operation

## Compatibility

| SDK versions        | Orchestrate versions       |
| ------------------- | -------------------------- |
| master/HEAD         | v21.1.0 (Unreleased)		   |
| Plugin v0.6.0       | v21.1.0 (Unreleased) 		   |

## Development

### Pre-requirements
- Go >= 1.14
- Makefile
- docker-compose

### Running local environment

To run our plugin in development mode you have to first build the plugin using:
```
$> make build
```

Then we run Hashicorp server with the following command:
```
$> vault server -dev -dev-plugin-dir=./build/bin/

...
    WARNING! dev mode is enabled! In this mode, Vault runs entirely in-memory
    and starts unsealed with a single unseal key. The root token is already
    authenticated to the CLI, so you can immediately begin using Vault.
    
    You may need to set the following environment variable:
    
        $ export VAULT_ADDR='http://127.0.0.1:8200'
    
    The unseal key and root token are displayed below in case you want to
    seal/unseal the Vault or re-authenticate.
    
    Unseal Key: 3hTanFX/q99PMBubXjwL/cFXh3YKCABPSmw31Jwok1w=
    Root Token: DevVaultToken
    
    The following dev plugins are registered in the catalog:
        - orchestrate-hashicorp-vault-plugin
```

Now you have your Vault running on port `:8200`. Open  a new terminal to run the following command to
enable Orchestrate plugin:
```
$> curl --header "X-Vault-Token: DevVaultToken" --request POST \
  --data '{"type": "plugin", "plugin_name": "orchestrate-hashicorp-vault-plugin", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/orchestrate
```

Now you already have your Vault running with Orchestrate plugin enable. The best way to understand the new
 integrate APIs is to use the `help` feature. To list a description of all the available endpoints you can run:
```
$> curl -H "X-Vault-Token: DevVaultToken" http://127.0.0.1:8200/v1/orchestrate?help=1
```

alternatively you can list only `ethereum` endpoints by using:
```
$> curl -H "X-Vault-Token: DevVaultToken" http://127.0.0.1:8200/v1/orchestrate/ethereum/accounts?help=1
```

#### Docker 
Alternatively, if you don't have HashiCorp installed locally you can run it using docker by:
```
$> make prod
```
To stop our environment we should run:
```
$> make down
```


## Production

Running Orchestrate plugin in production is slightly more complicated. First, we need to build plugin binary:
```
$> make build
```

Then copy our Orchestrate plugin to its own plugin folder, for instance, `/vault/plugins`
```
$> cp ./build/bin/orchestrate-hashicorp-vault-plugin /vault/plugins/orchestrate
```


We are setting our vault settings into a `config.hcl` where we indicate our plugins folder and rest of our
settings. You can find an example [here](./docker/config.hcl)

Now we initialize Hashicorp server using above settings
```
$> vault server -config=/vault/config.hcl
...
==> Vault server configuration:

             Api Address: http://vault:8200
                     Cgo: disabled
         Cluster Address: https://vault:8201
              Go Version: go1.15.4
              Listener 1: tcp (addr: "vault:8200", cluster address: "172.26.0.2:8201", max_request_duration: "1m30s", max_request_size: "33554432", tls: "disabled")
               Log Level: debug
                   Mlock: supported: true, enabled: true
           Recovery Mode: false
                 Storage: file
                 Version: Vault v1.6.1
             Version Sha: 6d2db3f033e02e70202bef9ec896360062b88b03

$> export VAULT_HOST=http://localhost:8200
```

Now we have to retrieve our environment UNSEAL_KEY and ROOT_TOKEN in order to manually enable our plugin
```
$> curl --request POST --data '{"secret_shares": 1, "secret_threshold": 1}' http://127.0.0.1/v1/sys/init
...
{"keys":["c04f4a5d21017a5ae4a421a083b251d65a0dc9ccdab3614e6bfbc6b452dc6c19"],"keys_base64":["wE9KXSEBelrkpCGgg7JR1loNyczas2FOa/vGtFLcbBk="],"root_token":"s.wI9JzbSooTwERi6aMrTE67kB"}

$> export UNSEAL_KEY=c04f4a5d21017a5ae4a421a083b251d65a0dc9ccdab3614e6bfbc6b452dc6c19
$> export ROOT_TOKEN= s.wI9JzbSooTwERi6aMrTE67kB
```

Unseal the Vault to allow changes on it:
```
$> curl --request POST --data '{"key": "${UNSEAL_KEY}"}' ${VAULT_HOST}/v1/sys/unseal
...
{"type":"shamir","initialized":true,"sealed":false,"t":1,"n":1,"progress":0,"nonce":"","version":"1.6.1","migration":false,"cluster_name":"vault-cluster-54f497d8","cluster_id":"b5f37577-bffe-e3ec-d2f3-338b5c7be4e5","recovery_seal":false,"storage_type":"file"}
```

For security reason to register a plugin into vault catalog we need to indicate its sha256 as follow:  
```
$> sha256sum -b /vault/plugin/orchestrate | cut -d' ' -f1
334227329451b9cea8dc4cc8c8e9700c206656b049c7020491795f7248b55029

$> curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data "{\"sha256\": \"334227329451b9cea8dc4cc8c8e9700c206656b049c7020491795f7248b55029\", \"command\": \"orchestrate\" }" \
  ${VAULT_HOST}/v1/sys/plugins/catalog/secret/orchestrate
```

Last step it would be to enable our plugin:
```
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "plugin", "plugin_name": "orchestrate", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_HOST}/v1/sys/mounts/orchestrate
```

### Docker
In case you are looking for a quick way to run a production ready using docker you have an example using:
```
$ make prod
```

## Orchestrate documentation

For a global understanding of Orchestrate, not only this Hashicorp Vault plugin, refer to the
[Orchestrate documentation.](https://docs.orchestrate.consensys.net/)


## Contributing
[How to Contribute](CONTRIBUTING.md)
