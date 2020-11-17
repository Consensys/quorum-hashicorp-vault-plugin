package main

import (
	"log"
	"os"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	client := &api.PluginAPIClientMeta{}
	err := client.FlagSet().Parse(os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: src.NewVaultBackend,
		TLSProviderFunc:    api.VaultPluginTLSProvider(client.GetTLSConfig()),
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
